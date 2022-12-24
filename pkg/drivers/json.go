package drivers

import (
	"encoding/json"
	"fmt"
	"github.com/SelfServiceCo/api/pkg/models"
	"os"
)

func JsonOpen(path string) *os.File {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	return jsonFile
}

func MysqlConfigLoad() models.Config {
	file, err := os.ReadFile("Config/conf.json")
	if err != nil {
		panic(err.Error())
	}

	var config models.Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		panic(err.Error())
	}

	return config
}
