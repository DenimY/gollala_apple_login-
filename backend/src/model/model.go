package model

type UserName struct {
	First string `json:"first"`
	Last  string `json:"last"`
	Full  string `json:"full"`
}
type SignInBody struct {
	Connect bool     `json:"connect"`
	UserId  string   `json:"userId"`
	Alias   string   `json:"alias"`
	Name    UserName `json:"name"`
}

type SignUpBody struct {
	Connect bool   `json:"connect"`
	UserId  string `json:"userId"`
}
