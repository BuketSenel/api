package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func GetMysqlConfig() (string, string, string) {
	vp := viper.New()
	vp.SetConfigName("mysql")
	vp.SetConfigType("json")
	vp.AddConfigPath("/www/configs/")
	err := vp.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	username := vp.GetString("username")
	database := vp.GetString("db")
	password := vp.GetString("password")
	return username, database, password
}
