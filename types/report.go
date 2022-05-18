package types

type Report struct {
	Id      string
	Type    string
	Address string
	Change  Diff
}

type Diff struct {
	Before string
	After  string
}
