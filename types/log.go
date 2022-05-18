package types

type RequestTrace struct {
	HttpMethod string
	StatusCode int64
	ID         string
	Content    string
}
