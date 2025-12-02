package controllers_model

import "time"

type LoginParams struct {
	Username string `form:"username" json:"username" binding:"required,alphanum"`
	Password string `form:"password" json:"password" binding:"required,min=6"`
	Token    string `form:"token_mobile" json:"token_mobile" binding:"required,min=6"`
}
type LogOutParams struct {
	RefreshToken string `form:"refresh_token" json:"refresh_token" binding:"required"`
}
type RegisterParams struct {
	Username string `form:"username" json:"username" binding:"required,alphanum"`
	Password string `form:"password" json:"password" binding:"required,min=6"`
	// Customer_info Customers `form:"customer" json:"customer" binding:"required"`
	Name   string    `form:"name" json:"name" binding:"required"`
	Email  string    `form:"email" json:"email" binding:"required,email"`
	Dob    time.Time `form:"dob" json:"dob" binding:"required"`
	Gender string    `form:"gender" json:"gender" binding:"required"`
}

type Customers struct {
	Name   string    `form:"name" json:"name"`
	Email  string    `form:"email" json:"email"`
	Dob    time.Time `form:"dob" json:"dob"`
	Gender string    `form:"gender" json:"gender"`
}

type CustomersAddress struct {
	Address_id  string `form:"address_id" json:"address_id"`
	Address     string `form:"address" json:"address"`
	PhoneNumber string `form:"phoneNumber" json:"phoneNumber" binding:"required,min=9,max=11"`
}
