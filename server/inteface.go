package server

// Validator is request validator interface
type Validator interface {
	Validate() error
}
