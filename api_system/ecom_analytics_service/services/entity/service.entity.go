package services

type Page struct {
	Page  int32 `json:"page" form:"page"`
	Limit int32 `json:"limit" form:"limit"`
	Total int   `json:"total" form:"-"`
}
type Narg[T any] struct {
	Data  T    `json:"data"`
	Valid bool `json:"valid"`
}
