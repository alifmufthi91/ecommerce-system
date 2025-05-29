package payload

type LoginUserReq struct {
	EmailOrPhone string `json:"email_or_phone" binding:"required"`
	Password     string `json:"password" binding:"required,min=8,max=32"`
}
