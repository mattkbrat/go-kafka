package data

type (
	UserType struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	UserParams struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	SignUpsProps struct {
		Id   string `json:"string"`
		User UserParams
	}
	SignInProps struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
)

type UsersType []UserParams

var (
	Users    = map[int]*UserParams{}
	SignUps  = map[int]*SignUpsProps{}
	Seq      = 0
	Sessions = map[string]*UserType{}
)
