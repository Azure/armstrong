package azidx

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/magodo/azure-rest-api-index/azidx/specpath"

	_ "embed"

	"github.com/go-git/go-git/v5"
	"github.com/go-openapi/jsonpointer"
	"github.com/go-openapi/jsonreference"
	"github.com/go-openapi/loads"
	"github.com/magodo/armid"
	"github.com/magodo/workerpool"
)

//go:embed dedup.json
var defaultDedup []byte

type FlattenOpIndex map[OpLocator]OperationRefs

type Index struct {
	Commit            string `json:"commit,omitempty"`
	ResourceProviders `json:"resource_providers"`
}

type ResourceProviders map[string]APIVersions

type APIVersions map[string]APIMethods

type APIMethods map[OperationKind]ResourceTypes

type ResourceTypes map[string]*OperationInfo

type OperationInfo struct {
	Actions       map[string]OperationRefs `json:"actions,omitempty"`
	OperationRefs OperationRefs            `json:"operation_refs,omitempty"`
}

const Wildcard = "*"

const ResourceRP = "Microsoft.Resources"

type OpLocator struct {
	// Upper cased RP name, e.g. MICROSOFT.COMPUTE. This might be "" for API path that has no explicit RP defined (e.g. /subscriptions/{subscriptionId})
	// This can be "*" to indicate it maps any RP
	RP string
	// API version, e.g. 2020-10-01-preview
	Version string
	// Upper cased resource type, e.g. /VIRTUALNETWORKS/SUBNETS
	// Each subtype can be "*" to indicate it maps any sub type.
	RT string
	// Upper cased potential action/collection type, e.g. LISTKEYS (action), SUBNETS
	// This can be "*" to indicate it maps any RP
	ACT string
	// HTTP operation kind, e.g. GET
	Method OperationKind
}

// OperationRefs represents a set of operation defintion (in form of JSON reference) that are mapped by the same operation locator.
// Since for a given operation locator, there might maps to multiple operation definition, only differing by the contained path pattern, there fore the actual operation ref is keyed by the containing path pattern.
// The value is a JSON reference to the operation, e.g. <dir>/foo.json#/paths/~1subscriptions~1{subscriptionId}~1providers~1{resourceProviderNamespace}~1register/post
type OperationRefs map[PathPatternStr]jsonreference.Ref

func (o OperationRefs) MarshalJSON() ([]byte, error) {
	m := map[string]string{}
	for k, v := range o {
		m[string(k)] = v.String()
	}
	return json.Marshal(m)
}

func (o *OperationRefs) UnmarshalJSON(b []byte) error {
	var m map[string]string
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	refs := OperationRefs{}
	for k, v := range m {
		v, err := url.PathUnescape(v)
		if err != nil {
			return err
		}
		refs[PathPatternStr(k)] = jsonreference.MustCreateRef(v)
	}
	*o = refs
	return nil
}

// PathPatternStr represents an API path pattern, with all the fixed segment upper cased, and all the parameterized segment as a literal "{}", or "{*}" (for x-ms-skip-url-encoding).
type PathPatternStr string

// BuildIndex builds the index file for the given specification directory.
// Since there are duplicated specification files in the directory, that defines the same API (same API path, version, operation), users can
// optionally specify a deduplication file. Otherwise, it will use a default dedup file instead.
func BuildIndex(specdir string, dedupFile string) (*Index, error) {
	specdir, err := filepath.Abs(specdir)
	if err != nil {
		return nil, err
	}

	b := defaultDedup
	if dedupFile != "" {
		b, err = os.ReadFile(dedupFile)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %v", dedupFile, err)
		}
	}
	var records DeduplicateRecords
	if err := json.Unmarshal(b, &records); err != nil {
		return nil, fmt.Errorf("unmarshal %s: %v", dedupFile, err)
	}
	deduplicator, err := records.ToDeduplicator()
	if err != nil {
		return nil, fmt.Errorf("converting the dedup file: %v", err)
	}

	var commit string
	repo, err := git.PlainOpen(filepath.Dir(specdir))
	if err != nil {
		if err != git.ErrRepositoryNotExists {
			return nil, err
		}
	} else {
		ref, err := repo.Head()
		if err != nil {
			return nil, err
		}
		commit = ref.Hash().String()
	}

	logger.Info("Collecting specs", "dir", specdir)
	l, err := collectSpecs(specdir)
	if err != nil {
		return nil, fmt.Errorf("collecting specs: %v", err)
	}
	logger.Info(fmt.Sprintf("%d specs collected", len(l)))

	logger.Info("Building operation index")
	ops, err := buildOpsIndex(specdir, deduplicator, l)
	if err != nil {
		return nil, fmt.Errorf("building operation index: %v", err)
	}

	// Turn flattend index to layerized index
	rps := ResourceProviders{}
	for loc, oprefs := range ops {
		rp, ok := rps[loc.RP]
		if !ok {
			rp = APIVersions{}
			rps[loc.RP] = rp
		}
		rpVer, ok := rp[loc.Version]
		if !ok {
			rpVer = APIMethods{}
			rp[loc.Version] = rpVer
		}
		rpVerMethod, ok := rpVer[loc.Method]
		if !ok {
			rpVerMethod = ResourceTypes{}
			rpVer[loc.Method] = rpVerMethod
		}
		rpVerMethodRt, ok := rpVerMethod[loc.RT]
		if !ok {
			rpVerMethodRt = &OperationInfo{}
			rpVerMethod[loc.RT] = rpVerMethodRt
		}
		if loc.ACT != "" {
			if rpVerMethodRt.Actions == nil {
				rpVerMethodRt.Actions = map[string]OperationRefs{}
			}
			if _, ok := rpVerMethodRt.Actions[loc.ACT]; ok {
				return nil, fmt.Errorf("unexpected duplicated resource action at %#v", loc)
			}
			rpVerMethodRt.Actions[loc.ACT] = oprefs
		} else {
			rpVerMethodRt.OperationRefs = oprefs
		}
	}

	index := &Index{
		Commit:            commit,
		ResourceProviders: rps,
	}

	return index, nil
}

// collectSpecs collects all Swagger specs based on the effective tags in each RP's readme.md.
func collectSpecs(rootdir string) ([]string, error) {
	var speclist []string

	if err := filepath.WalkDir(rootdir,
		func(p string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if strings.EqualFold(d.Name(), "data-plane") {
					return filepath.SkipDir
				}
				if strings.EqualFold(d.Name(), "examples") {
					return filepath.SkipDir
				}
				return nil
			}
			if d.Name() != "readme.md" {
				return nil
			}
			content, err := os.ReadFile(p)
			if err != nil {
				return fmt.Errorf("reading file %s: %v", p, err)
			}
			l, err := SpecListFromReadmeMD(content)
			if err != nil {
				return fmt.Errorf("retrieving spec list from %s: %v", p, err)
			}
			for _, relp := range l {
				speclist = append(speclist, filepath.Join(filepath.Dir(p), relp))
			}
			return filepath.SkipDir
		}); err != nil {
		return nil, err
	}
	// Deduplicate
	m := map[string]struct{}{}
	for _, v := range speclist {
		m[v] = struct{}{}
	}
	speclist = make([]string, 0, len(m))
	for k := range m {
		speclist = append(speclist, k)
	}
	// Sort
	sort.Slice(speclist, func(i, j int) bool { return speclist[i] < speclist[j] })

	return speclist, nil
}

func buildOpsIndex(specdir string, deduplicator Deduplicator, specs []string) (FlattenOpIndex, error) {
	specdir, err := filepath.Abs(specdir)
	if err != nil {
		return nil, err
	}
	ops := FlattenOpIndex{}
	var lock sync.Mutex

	type dupkey struct {
		OpLocator
		PathPatternStr
	}
	dups := map[dupkey][]jsonreference.Ref{}

	wp := workerpool.NewWorkPool(runtime.NumCPU())
	wp.Run(nil)
	for _, spec := range specs {
		spec := spec
		wp.AddTask(func() (interface{}, error) {
			m, err := parseSpec(specdir, spec)
			if err != nil {
				return nil, fmt.Errorf("parsing spec %s: %v", spec, err)
			}

			lock.Lock()
			defer lock.Unlock()

			for k, mm := range m {
				if len(ops[k]) == 0 {
					ops[k] = OperationRefs{}
				}
				for ppattern, ref := range mm {
					if exist, ok := ops[k][ppattern]; ok {
						// Temporarily record duplicate operation definitions and resolve it later
						k := dupkey{
							OpLocator:      k,
							PathPatternStr: ppattern,
						}
						if len(dups[k]) == 0 {
							dups[k] = append(dups[k], exist)
						}
						dups[k] = append(dups[k], ref)
						continue
					}
					ops[k][ppattern] = ref
				}
			}
			return nil, nil
		})
	}
	if err := wp.Done(); err != nil {
		return nil, err
	}

	// Resolve duplicates (auto)
	newdups := map[dupkey][]jsonreference.Ref{}
	for k, refs := range dups {
		newrefs := make([]jsonreference.Ref, len(refs))
		copy(newrefs, refs)
		newdups[k] = newrefs
	}
	for k, refs := range dups {
		var candidateRefs []jsonreference.Ref
		for _, ref := range refs {
			pinfo, err := specpath.SpecPathInfo(ref.GetURL().Path)
			if err != nil {
				return nil, fmt.Errorf("new spec path info: %v", err)
			}
			// Only pick up the op locator that well matches its spec path, which hopefully is the orignal spec that defines this operation
			if strings.EqualFold(pinfo.ResourceProviderMS, k.RP) &&
				strings.EqualFold(pinfo.Version, k.Version) &&
				pinfo.IsPreview == strings.HasSuffix(k.Version, "preview") {
				candidateRefs = append(candidateRefs, ref)
			}
		}
		if len(candidateRefs) == 1 {
			ops[k.OpLocator][k.PathPatternStr] = candidateRefs[0]
			delete(newdups, k)
		}
	}
	dups = newdups

	// Resolve duplicates (manually)
	for k, refs := range dups {
		var dedupOp *DedupOp
		var matcherName string

		// Look for the dedup operator
		if deduplicator != nil {
			for matcher, op := range deduplicator {
				op := op
				if matcher.Match(k.OpLocator, string(k.PathPatternStr)) {
					if dedupOp != nil {
						return nil, fmt.Errorf("Duplicate matchers in duplicator that match %s: %s vs %s", k, matcherName, matcher.Name)
					}
					dedupOp = &op
					matcherName = matcher.Name
				}
			}
		}

		var refStrs []string
		for _, ref := range refs {
			refStrs = append(refStrs, ref.String())
		}
		refMsg := "\n" + strings.Join(refStrs, "\n")

		if dedupOp != nil {
			switch {
			case dedupOp.Picker != nil:
				picker := dedupOp.Picker
				var pickCnt int
				var pickRef jsonreference.Ref
				for _, ref := range refs {
					if picker.Match(ref) {
						pickCnt++
						pickRef = ref
					}
				}
				if pickCnt == 0 {
					logger.Warn("dedup matcher picked nothing", "oploc", k.OpLocator, "path", k.PathPatternStr, "matcher", matcherName, "refs", refMsg)
					continue
					//return nil, fmt.Errorf("dedup matcher %s picked nothing for %s. refs: %v", matcherName, k, refMsg)
				}
				if pickCnt > 1 {
					logger.Warn("still have duplicates after dedup picking", "oploc", k.OpLocator, "path", k.PathPatternStr, "matcher", matcherName, "refs", refMsg)
					continue
					//return nil, fmt.Errorf("still have duplicates after dedup matcher %s picking for %s. refs: %v", matcherName, k, refMsg)
				}
				ops[k.OpLocator][k.PathPatternStr] = pickRef
			case dedupOp.Any:
				ops[k.OpLocator][k.PathPatternStr] = refs[0]
			case dedupOp.Ignore:
				delete(ops[k.OpLocator], k.PathPatternStr)
				if len(ops[k.OpLocator]) == 0 {
					delete(ops, k.OpLocator)
				}
			}
			continue
		}

		logger.Warn("duplicate definition", "oploc", k.OpLocator, "path", k.PathPatternStr, "refs", refMsg)
	}

	return ops, nil
}

// parseSpec parses one Swagger spec and returns back a operation index for this spec
func parseSpec(specdir, p string) (FlattenOpIndex, error) {
	doc, err := loads.Spec(p)
	if err != nil {
		return nil, fmt.Errorf("loading spec: %v", err)
	}
	swagger := doc.Spec()

	// Skipping swagger specs that have no "paths" defined
	if swagger.Paths == nil || len(swagger.Paths.Paths) == 0 {
		return nil, nil
	}
	if swagger.Info == nil {
		return nil, fmt.Errorf(`spec has no "Info"`)
	}
	if swagger.Info.Version == "" {
		return nil, fmt.Errorf(`spec has no "Info.Version"`)
	}

	absSpecPath, err := filepath.Abs(p)
	if err != nil {
		return nil, fmt.Errorf("failed to get abs path for %s: %v", p, err)
	}
	relSpecPath, err := filepath.Rel(specdir, absSpecPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get rel path for %s: %v", p, err)
	}

	pinfo, err := specpath.SpecPathInfo(relSpecPath)
	if err != nil {
		return nil, fmt.Errorf("new spec path info: %v", err)
	}

	version := swagger.Info.Version
	index := FlattenOpIndex{}
	for path, pathItem := range swagger.Paths.Paths {
		for _, opKind := range PossibleOperationKinds {
			if PathItemOperation(pathItem, opKind) == nil {
				continue
			}
			logger.Debug("Parsing spec", "spec", p, "path", path, "operation", opKind)
			pathPatterns, err := ParsePathPatternFromSwagger(p, swagger, path, opKind)
			if err != nil {
				return nil, fmt.Errorf("parsing path pattern for %s (%s): %v", path, opKind, err)
			}
			for _, pathPattern := range pathPatterns {
				// path -> RP, RT, ACT
				// We look backwards for the first "providers" segment.
				providerIdx := -1
				for i := len(pathPattern.Segments) - 1; i >= 0; i-- {
					if strings.EqualFold(pathPattern.Segments[i].FixedName, "providers") {
						providerIdx = i
						break
					}
				}
				var (
					rp, rt, act string
					rpIsGlob    bool
					nextIdx     int
				)

				if providerIdx == -1 || len(pathPattern.Segments) == providerIdx+1 {
					// No "providers" segment found, if the spec is from Resources RP, then add up the implicit RP "Microsoft.Resources"
					if !strings.EqualFold(pinfo.ResourceProviderMS, "Microsoft.Resources") {
						logger.Warn("no provider defined", "spec", p, "path", path, "operation", opKind, "rp", pinfo.ResourceProviderMS)
						continue
					}
					rp = ResourceRP
					nextIdx = 0
				} else {
					// RP found
					providerSeg := pathPattern.Segments[providerIdx+1]
					rp = providerSeg.FixedName
					rpIsGlob = providerSeg.IsParameter
					nextIdx = providerIdx + 2
				}

				// Ignore the too generic api paths:
				// 1. Those have only one multi-segmented parameter segment. E.g. /{resourceId}
				if len(pathPattern.Segments) == 1 {
					continue
				}
				// 2. Those whose provider and all the following segments are parameterized. E.g. /subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{parentResourcePath}/{resourceType}/{resourceName}
				if rpIsGlob && allParameterized(pathPattern.Segments[nextIdx:]) {
					continue
				}

				// Identify the ACT
				lastIdx := len(pathPattern.Segments)
				if len(pathPattern.Segments[nextIdx:])%2 == 1 {
					lastIdx = lastIdx - 1
					seg := pathPattern.Segments[len(pathPattern.Segments)-1]
					if seg.IsParameter {
						if seg.IsMulti {
							logger.Warn("action segment is multi-segmented parameter", "path", path, "operation", opKind)
							continue
						}
						act = Wildcard
					} else {
						act = seg.FixedName
					}
				}

				// Identify the RT
				var rts []string
				for i := nextIdx; i < lastIdx; i += 2 {
					seg := pathPattern.Segments[i]
					var rtName string
					if seg.IsParameter {
						if seg.IsMulti {
							logger.Warn("resource type is multi-segmented parameter", "path", path, "operation", opKind, "index", i)
							continue
						}
						rtName = Wildcard
					} else {
						rtName = seg.FixedName
					}
					rts = append(rts, rtName)
				}
				rt = "/" + strings.Join(rts, "/")

				opRef := jsonreference.MustCreateRef(relSpecPath + "#/paths/" + jsonpointer.Escape(path) + "/" + strings.ToLower(string(opKind)))

				pathPatternStr := PathPatternStr(strings.ToUpper(pathPattern.String()))

				opLoc := OpLocator{
					RP:      strings.ToUpper(rp),
					Version: version,
					RT:      strings.ToUpper(rt),
					ACT:     strings.ToUpper(act),
					Method:  opKind,
				}

				if rpIsGlob {
					opLoc.RP = Wildcard
				}

				if _, ok := index[opLoc]; !ok {
					index[opLoc] = map[PathPatternStr]jsonreference.Ref{}
				}
				if exist, ok := index[opLoc][pathPatternStr]; ok {
					return nil, fmt.Errorf(
						"operation locator %#v for path pattern %s already applied with operation %s, conflicts to the new operation %s", opLoc, pathPatternStr, &exist, &opRef)
				}
				index[opLoc][pathPatternStr] = opRef
			}
		}
	}
	return index, nil
}

func (idx Index) Lookup(method string, uRL url.URL) (*jsonreference.Ref, error) {
	operation := OperationKind(strings.ToUpper(method))
	apiVersion := uRL.Query().Get("api-version")

	path := strings.TrimRight(strings.ToUpper(uRL.Path), "/")
	segs := strings.Split(strings.TrimLeft(path, "/"), "/")

	respath := path
	var act string
	if len(segs)%2 == 1 {
		act = strings.ToUpper(segs[len(segs)-1])
		respath = "/" + strings.Join(segs[:len(segs)-1], "/")
	}
	id, err := armid.ParseResourceId(respath)
	if err != nil {
		return nil, fmt.Errorf("parsing %s as arm id: %v", respath, err)
	}

	rp := strings.ToUpper(id.Provider())
	rt := strings.ToUpper("/" + strings.Join(id.Types(), "/"))

	if rpInfo, ok := idx.ResourceProviders[rp]; ok {
		ref, ok, err := lookupIntoRP(rpInfo, path, apiVersion, operation, rt, act)
		if err != nil {
			return nil, fmt.Errorf("lookup for %v (%s) in rp %s: %v", uRL.String(), method, rp, err)
		}
		if ok {
			return ref, nil
		}
	}
	ref, ok, err := lookupIntoRP(idx.ResourceProviders[Wildcard], path, apiVersion, operation, rt, act)
	if err != nil {
		return nil, fmt.Errorf("lookup for %v (%s) in the wildcard rp: %v", uRL.String(), method, err)
	}
	if !ok {
		return nil, fmt.Errorf("lookup for %v (%s): matches nothing", uRL.String(), method)
	}
	return ref, nil
}

func lookupIntoRP(rpInfo map[string]APIMethods, path, apiVersion string, operation OperationKind, rt, act string) (*jsonreference.Ref, bool, error) {
	rpVer, ok := rpInfo[apiVersion]
	if !ok {
		return nil, false, nil
	}
	rpVerOp, ok := rpVer[operation]
	if !ok {
		return nil, false, nil
	}

	buildRTMatcher := func(rt string) Matcher {
		segs := strings.Split(strings.Trim(rt, "/"), "/")
		m := Matcher{
			PrefixSep: true,
			Separater: "/",
		}
		for _, seg := range segs {
			if seg == "*" {
				m.Segments = append(m.Segments, MatchSegment{IsWildcard: true})
				continue
			}
			m.Segments = append(m.Segments, MatchSegment{Value: seg})
		}
		return m
	}

	type opInfoWrapper struct {
		info      OperationInfo
		rtMatcher Matcher
	}

	var opInfoWrappers []opInfoWrapper
	for rt, opInfo := range rpVerOp {
		opInfoWrappers = append(opInfoWrappers, opInfoWrapper{
			info:      *opInfo,
			rtMatcher: buildRTMatcher(rt),
		})
	}

	// Sort the resource type matchers to match from the most specific to the most general
	sort.Slice(opInfoWrappers, func(i, j int) bool {
		return opInfoWrappers[i].rtMatcher.Less(opInfoWrappers[j].rtMatcher)
	})

	for _, opInfoWrapper := range opInfoWrappers {
		if !opInfoWrapper.rtMatcher.Match(rt) {
			continue
		}
		opInfo := opInfoWrapper.info
		oprefs := opInfo.OperationRefs
		if act != "" {
			if len(opInfo.Actions) == 0 {
				continue
			}
			var ok bool
			oprefs, ok = opInfo.Actions[act]
			if !ok {
				oprefs, ok = opInfo.Actions[Wildcard]
				if !ok {
					continue
				}
			}
		}

		// Select the best matching path from candidate paths
		type opRefWrapper struct {
			ref         jsonreference.Ref
			pathMatcher Matcher
		}
		var opRefWrappers []opRefWrapper
		for ppath, ref := range oprefs {
			pathPattern := ParsePathPatternFromString(string(ppath))
			m := Matcher{
				PrefixSep: true,
				Separater: "/",
			}
			for _, seg := range pathPattern.Segments {
				m.Segments = append(m.Segments, MatchSegment{
					Value:      seg.FixedName,
					IsWildcard: seg.IsParameter,
					IsAny:      seg.IsMulti,
				})
			}
			opRefWrappers = append(opRefWrappers, opRefWrapper{
				ref:         ref,
				pathMatcher: m,
			})
		}

		// Sort the path matchers to match from the most specific to the most general
		sort.Slice(opRefWrappers, func(i, j int) bool {
			return opRefWrappers[i].pathMatcher.Less(opRefWrappers[j].pathMatcher)
		})

		for _, opRefWrapper := range opRefWrappers {
			if opRefWrapper.pathMatcher.Match(path) {
				return &opRefWrapper.ref, true, nil
			}
		}
	}

	return nil, false, nil
}

func allParameterized(segs []PathSegment) bool {
	for _, seg := range segs {
		if !seg.IsParameter {
			return false
		}
	}
	return true
}
