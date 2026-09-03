package errors

type RequestError struct {
	Err     error
	Message string
}
