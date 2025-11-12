package controllers_model

type RatingParams struct {
	ProductsSpuID string `form:"products_spu_id" json:"products_spu_id" binding:"required"`
	Star          int    `form:"star" json:"star" binding:"required"`
	Comment       string `form:"comment" json:"comment" `
}
