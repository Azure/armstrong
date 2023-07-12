package specpath

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Info struct {
	ResourceProvider   string
	ResourceProviderMS string
	IsPreview          bool
	Version            string
	SpecName           string

	// subservice is used to store the optional sub-servce after the RP if any
	subservice *string
}

// SpecPathInfo returns the SpecPathInfo of the Path, given the Path is a relative path to a swagger spec.
// e.g.
// - compute/resource-manager/Microsoft.Compute/stable/2020-01-01/compute.json
// - mediaservices/resource-manager/Microsoft.Media/Accounts/preview/2019-05-01-preview/Accounts.json
func SpecPathInfo(p string) (*Info, error) {
	if filepath.IsAbs(p) {
		return nil, fmt.Errorf("expect relative path to the spec from spec rootdir, got %s", p)
	}
	segs := strings.Split(string(p), string(os.PathSeparator))
	if len(segs) < 6 {
		return nil, fmt.Errorf("swagger spec path expects more nesting level, but %q has only level %d", p, len(segs))
	}
	if len(segs) > 7 {
		return nil, fmt.Errorf("swagger spec path expects less nesting level, but %q has only level %d", p, len(segs))
	}
	filename := filepath.Base(string(p))
	if filepath.Ext(filename) != ".json" {
		return nil, fmt.Errorf("expect a json file, got %s", filename)
	}

	var subservice *string
	if len(segs) == 7 {
		subservice = &segs[len(segs)-4]
	}

	return &Info{
		ResourceProvider:   segs[0],
		ResourceProviderMS: segs[2],
		IsPreview:          segs[len(segs)-3] == "preview",
		Version:            segs[len(segs)-2],
		SpecName:           filename,

		subservice: subservice,
	}, nil
}

// ToPath returns the relative path of the spec file
func (info Info) ToPath() string {
	segs := []string{info.ResourceProvider, "resource-manager", info.ResourceProviderMS}
	if info.subservice != nil {
		segs = append(segs, *info.subservice)
	}
	if info.IsPreview {
		segs = append(segs, "preview")
	} else {
		segs = append(segs, "stable")
	}
	segs = append(segs, info.Version, info.SpecName)
	return filepath.Join(segs...)
}
