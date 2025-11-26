package assets_services

import "fmt"

type ServiceError struct {
	Code int
	Err  error
}

func NewError(code int, err error) *ServiceError {
	return &ServiceError{Code: code, Err: err}
}
func (e *ServiceError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return fmt.Sprintf("ServiceError with code: %d", e.Code)
}
