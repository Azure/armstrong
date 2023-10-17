package types

import paltypes "github.com/ms-henglu/pal/types"

type PassReport struct {
	Resources []Resource
}

type Resource struct {
	Type    string
	Address string
}

type DiffReport struct {
	Diffs []Diff
	Logs  []paltypes.RequestTrace
}

type Diff struct {
	Id      string
	Type    string
	Address string
	Change  Change
}

type Change struct {
	Before string
	After  string
}

type ErrorReport struct {
	Errors []Error
	Logs   []paltypes.RequestTrace
}

type Error struct {
	Id      string
	Type    string
	Label   string
	Message string
}
