package helper

import (
	"github.com/ms-henglu/armstrong/types"
	"github.com/nsf/jsondiff"
)

func DiffMessageTerraform(diff types.Diff) string {
	option := jsondiff.DefaultConsoleOptions()
	_, msg := jsondiff.Compare([]byte(diff.Before), []byte(diff.After), &option)
	return msg
}

func DiffMessageReadable(diff types.Diff) string {
	option := jsondiff.Options{
		Added:                 jsondiff.Tag{Begin: "\033[0;32m", End: " is not returned from response\033[0m"},
		Removed:               jsondiff.Tag{Begin: "\033[0;31m", End: "\033[0m"},
		Changed:               jsondiff.Tag{Begin: "\033[0;33m Got ", End: "\033[0m"},
		Skipped:               jsondiff.Tag{Begin: "\033[0;90m", End: "\033[0m"},
		SkippedArrayElement:   jsondiff.SkippedArrayElement,
		SkippedObjectProperty: jsondiff.SkippedObjectProperty,
		ChangedSeparator:      " in response, expect ",
		Indent:                "    ",
	}
	_, msg := jsondiff.Compare([]byte(diff.Before), []byte(diff.After), &option)
	return msg
}

func DiffMessageMarkdown(diff types.Diff) string {
	option := jsondiff.Options{
		Added:                 jsondiff.Tag{Begin: "", End: " is not returned from response"},
		Removed:               jsondiff.Tag{Begin: "", End: ""},
		Changed:               jsondiff.Tag{Begin: "Got ", End: ""},
		Skipped:               jsondiff.Tag{Begin: "", End: ""},
		SkippedArrayElement:   jsondiff.SkippedArrayElement,
		SkippedObjectProperty: jsondiff.SkippedObjectProperty,
		ChangedSeparator:      " in response, expect ",
		Indent:                "    ",
	}
	_, msg := jsondiff.Compare([]byte(diff.Before), []byte(diff.After), &option)
	return msg
}
