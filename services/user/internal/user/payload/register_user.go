package payload

type RegisterUserReq struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required,min=3,max=50"`
	Phone    string `json:"phone" binding:"required,min=10,max=15"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}
