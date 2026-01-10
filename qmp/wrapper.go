package qmp

type Wrapper[T any] struct {
	Type string `json:"type"`
	Data T      `json:"data"`
}

type Response[T any] struct {
	Return T `json:"return"`
}
