package azidx

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-openapi/spec"
)

type PathPattern struct {
	Segments []PathSegment
}

type PathSegment struct {
	FixedName   string
	IsParameter bool
	IsMulti     bool // indicates the x-ms-skip-url-encoding = true
}

func ParsePathPatternFromSwagger(specFile string, swagger *spec.Swagger, path string, operation OperationKind) ([]PathPattern, error) {
	if swagger.Paths == nil {
		return nil, fmt.Errorf(`no "paths"`)
	}
	pathItem, ok := swagger.Paths.Paths[path]
	if !ok {
		return nil, fmt.Errorf(`no path %s found`, path)
	}
	parameterMap := map[string]spec.Parameter{}
	for _, param := range pathItem.Parameters {
		if param.Ref.String() != "" {
			pparam, err := spec.ResolveParameterWithBase(swagger, param.Ref, &spec.ExpandOptions{RelativeBase: specFile})
			if err != nil {
				return nil, fmt.Errorf("resolving ref %q: %v", param.Ref.String(), err)
			}
			param = *pparam
		}
		parameterMap[param.Name] = param
	}
	// Per operation parameter overrides the per path parameter
	if op := PathItemOperation(pathItem, operation); op != nil {
		for _, param := range op.Parameters {
			if param.Ref.String() != "" {
				pparam, err := spec.ResolveParameterWithBase(swagger, param.Ref, &spec.ExpandOptions{RelativeBase: specFile})
				if err != nil {
					return nil, fmt.Errorf("resolving ref %q: %v", param.Ref.String(), err)
				}
				param = *pparam
			}
			parameterMap[param.Name] = param
		}
	}

	// Initialliy, there is only one []PathSegment in the segment set.
	segmentSet := [][]PathSegment{{}}
	addSegment := func(sset [][]PathSegment, s PathSegment) {
		for i := 0; i < len(sset); i++ {
			sset[i] = append(sset[i], s)
		}
	}
	for _, seg := range strings.Split(strings.Trim(path, "/"), "/") {
		if isParameterizedSegment(seg) {
			segment := PathSegment{
				IsParameter: true,
			}

			// There are very limited API paths that define more than one parameters in one segment, e.g.:
			// https://github.com/Azure/azure-rest-api-specs/blob/b672a0b301338a570af2e5430b4b7691f909a094/specification/eventgrid/resource-manager/Microsoft.EventGrid/preview/2023-12-15-preview/EventGrid.json#L9098
			// In this case, we simply treat it as a single parameter for now, to not complicate things too much for a minority cases
			// (actually, currently there is only such one, and doing so won't affect the index lookup result at all).
			if hasMultipleParameterizedSegment(seg) {
				addSegment(segmentSet, segment)
				continue
			}

			// In case this segment is an enum parameter, replicate the existing patterns to time of (the amount of enum variants with the variant + 1 of the original wildcard) appended
			name := strings.Trim(seg, "{}")
			param, ok := parameterMap[name]
			if !ok {
				return nil, fmt.Errorf("undefined parameter name %q", name)
			}
			if param.HasEnum() {
				var newSegmentSet [][]PathSegment
				for _, enum := range param.Enum {
					for _, segs := range segmentSet {
						newSegs := make([]PathSegment, len(segs)+1)
						copy(newSegs, segs)
						// enum parameter in path must be of type string
						newSegs[len(newSegs)-1] = PathSegment{FixedName: enum.(string)}
						newSegmentSet = append(newSegmentSet, newSegs)
					}
				}
				for _, segs := range segmentSet {
					newSegs := make([]PathSegment, len(segs)+1)
					copy(newSegs, segs)
					newSegs[len(newSegs)-1] = PathSegment{IsParameter: true}
					newSegmentSet = append(newSegmentSet, newSegs)
				}
				segmentSet = newSegmentSet
				continue
			}

			// Skip URL encoding
			if v, ok := param.VendorExtensible.Extensions["x-ms-skip-url-encoding"]; ok && v.(bool) {
				segment.IsMulti = true
			}

			addSegment(segmentSet, segment)
			continue
		}

		// Non parameterized segment
		addSegment(segmentSet, PathSegment{FixedName: seg})
	}

	var pathPatterns []PathPattern
	for _, segs := range segmentSet {
		pathPatterns = append(pathPatterns, PathPattern{Segments: segs})
	}
	sort.Slice(pathPatterns, func(i, j int) bool { return pathPatterns[i].String() < pathPatterns[j].String() })
	return pathPatterns, nil
}

func ParsePathPatternFromString(path string) *PathPattern {
	var segments []PathSegment
	for _, seg := range strings.Split(strings.Trim(path, "/"), "/") {
		switch seg {
		case "{}":
			segments = append(segments, PathSegment{IsParameter: true})
		case "{*}":
			segments = append(segments, PathSegment{IsParameter: true, IsMulti: true})
		default:
			segments = append(segments, PathSegment{FixedName: seg})
		}
	}
	return &PathPattern{Segments: segments}
}

func isParameterizedSegment(seg string) bool {
	return strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}")
}

func hasMultipleParameterizedSegment(seg string) bool {
	var (
		found bool
		left  bool
	)
	for _, c := range seg {
		if c != '{' && c != '}' {
			continue
		}
		if c == '{' {
			left = true
			continue
		}
		// this is '}'
		if left {
			if found {
				return true
			}
			found = true
			left = false
			continue
		}
	}
	return false
}

func (p PathPattern) String() string {
	var segs []string
	for _, seg := range p.Segments {
		if !seg.IsParameter {
			segs = append(segs, seg.FixedName)
			continue
		}
		if seg.IsMulti {
			segs = append(segs, "{*}")
		} else {
			segs = append(segs, "{}")
		}
	}
	return "/" + strings.Join(segs, "/")
}
