package models

type Config struct {
	Name     string `json:"name"`
	Db       string `json:"db"`
	Password string `json:"password"`
}
