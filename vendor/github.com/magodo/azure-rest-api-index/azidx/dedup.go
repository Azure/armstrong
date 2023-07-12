package azidx

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-openapi/jsonreference"
)

type DedupMatcher struct {
	Name        string
	RP          *regexp.Regexp
	Version     *regexp.Regexp
	RT          *regexp.Regexp
	ACT         *regexp.Regexp
	Method      *regexp.Regexp
	PathPattern *regexp.Regexp
}

func (key DedupMatcher) Match(loc OpLocator, pathp string) bool {
	return (key.RP == nil || key.RP.MatchString(loc.RP)) &&
		(key.Version == nil || key.Version.MatchString(loc.Version)) &&
		(key.RT == nil || key.RT.MatchString(loc.RT)) &&
		(key.ACT == nil || key.ACT.MatchString(loc.ACT)) &&
		(key.Method == nil || key.Method.MatchString(string(loc.Method))) &&
		(key.PathPattern == nil || key.PathPattern.MatchString(pathp))
}

type DedupPicker struct {
	SpecPath *regexp.Regexp
	Pointer  *regexp.Regexp
}

func (picker DedupPicker) Match(ref jsonreference.Ref) bool {
	specPath := ref.GetURL().Path
	pointer := ref.GetPointer().String()
	return (picker.SpecPath == nil || picker.SpecPath.MatchString(specPath)) &&
		(picker.Pointer == nil || picker.Pointer.MatchString(pointer))
}

type Deduplicator map[DedupMatcher]DedupOp

type DedupOp struct {
	Picker *DedupPicker
	Ignore bool
	Any    bool
}

type DeduplicateRecords map[string]DedupRecord

type DedupRecord struct {
	Matcher DedupMatcherIn `json:"matcher"`
	Picker  *DedupPickerIn `json:"picker"`
	Ignore  *bool          `json:"ignore"`
	Any     *bool          `json:"any"`
}

type DedupMatcherIn struct {
	RP      string   `json:"rp,omitempty"`
	Version string   `json:"version,omitemptyn"`
	RT      string   `json:"rt,omitempty"`
	ACT     string   `json:"act,omitempty"`
	Method  string   `json:"method,omitempty"`
	Paths   []string `json:"paths,omitempty"`
}

type DedupPickerIn struct {
	SpecPath string `json:"spec_path,omitempty"`
	Pointer  string `json:"pointer,omitempty"`
}

func (records DeduplicateRecords) ToDeduplicator() (Deduplicator, error) {
	dup := Deduplicator{}
	for name, rec := range records {
		matcher := rec.Matcher
		m := DedupMatcher{Name: name}
		if matcher.RP != "" {
			m.RP = regexp.MustCompile(matcher.RP)
		}
		if matcher.Version != "" {
			m.Version = regexp.MustCompile(matcher.Version)
		}
		if matcher.RT != "" {
			m.RT = regexp.MustCompile(matcher.RT)
		}
		if matcher.ACT != "" {
			m.ACT = regexp.MustCompile(matcher.ACT)
		}
		if matcher.Method != "" {
			m.Method = regexp.MustCompile(matcher.Method)
		}
		if len(matcher.Paths) != 0 {
			var pstrs []string
			for _, p := range matcher.Paths {
				pstrs = append(pstrs, "^"+p+"$")
			}
			m.PathPattern = regexp.MustCompile(strings.Join(pstrs, "|"))
		}
		op := DedupOp{}
		var n int
		if rec.Picker != nil {
			n++
		}
		if rec.Ignore != nil && *rec.Ignore {
			n++
		}
		if rec.Any != nil && *rec.Any {
			n++
		}
		if n != 1 {
			return nil, fmt.Errorf("exactly one of `ignore`, `any` and `picker` has to be specified")
		}
		if rec.Ignore != nil {
			op.Ignore = *rec.Ignore
		}
		if rec.Any != nil {
			op.Any = *rec.Any
		}
		if rec.Picker != nil {
			picker := rec.Picker
			p := DedupPicker{}
			if picker.SpecPath != "" {
				p.SpecPath = regexp.MustCompile(picker.SpecPath)
			}
			if picker.Pointer != "" {
				p.Pointer = regexp.MustCompile(picker.Pointer)
			}
			op.Picker = &p
		}
		dup[m] = op
	}
	return dup, nil
}
