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

// ProductRating represents rating information for a product
type ProductRating struct {
	ProductID     string  `json:"product_id"`
	TotalReviews  int64   `json:"total_reviews"`
	AverageRating float64 `json:"average_rating"`
}
