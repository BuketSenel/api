package models

type Config struct {
	Name          string `json:"name"`
	Db            string `json:"db"`
	Password      string `json:"password"`
	BUCKET_REGION string `json:"S3_REGION"`
	BUCKET_NAME   string `json:"S3_BUCKET"`
	BUCKET_SECRET string `json:"S3_SECRET"`
	BUCKET_KEY    string `json:"S3_KEY"`
}
