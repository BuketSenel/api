package models

type Config struct {
	Name          string `json:"name"`
	Db            string `json:"db"`
	Password      string `json:"password"`
	BUCKET_REGION string `json:"BUCKET_REGION"`
	BUCKET_NAME   string `json:"BUCKET_NAME"`
	BUCKET_SECRET string `json:"BUCKET_SECRET"`
	BUCKET_KEY    string `json:"BUCKET_KEY"`
}
