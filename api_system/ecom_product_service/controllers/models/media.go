package controllers_model

type DeleteMediaParams struct {
	ListID []string `form:"list_url" json:"list_url" binding:"required"`
}
