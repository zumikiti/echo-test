package user

type Request struct {
	ID    int64
	Name  string `json:"name" form:"name"`
	Email string `json:"email" form:"email"`
}
