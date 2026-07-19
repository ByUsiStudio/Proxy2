package bean

type ReqRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Desc     string `json:"desc"`
}
