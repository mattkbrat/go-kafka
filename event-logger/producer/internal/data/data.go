package data

type UserType struct {
	Name     string `json:"name"`
	Username string `json:"Username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpsType struct {
	Id   string `json:"string"`
	User UserType
}

type UsersType []UserType

var (
	Users   = map[int]*UserType{}
	SignUps = map[int]*SignUpsType{}
	Seq     = 0
)
