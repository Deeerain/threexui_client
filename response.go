package threexuiclient

type XUIResponse[T any] struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Obj     T      `json:"obj"`
}
