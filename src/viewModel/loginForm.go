package viewModel

type LoginForm struct {
	Username *string `json:"username" form:"username" binding:"required"`
	Password *string `json:"password" from:"password" binding:"required"`
}
