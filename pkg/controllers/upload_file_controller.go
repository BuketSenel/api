package controllers

import (
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) (string, gin.H) {
	file, err := c.FormFile("file")
	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": "Could not retrive the file from the request!"}
	}
	rest_id := c.PostForm("rest_id")

	filePath := "/www/uploads/" + rest_id + "/" + file.Filename
	err = c.SaveUploadedFile(file, filePath)

	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": err.Error()}
	}

	f, err := os.Open(filePath)
	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": err.Error()}
	}
	session, err := session.NewSession(&aws.Config{Region: aws.String(conf.BUCKET_REGION), Credentials: credentials.NewStaticCredentials(conf.BUCKET_KEY, conf.BUCKET_SECRET, "")})
	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": "Error creating session!"}
	}

	_, err = s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(conf.BUCKET_NAME),
		Key:    aws.String(rest_id + "/" + file.Filename),
		ACL:    aws.String("private"),
		Body:   f,
	})
	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": err.Error()}
	}
	domain := "https://s3." + conf.BUCKET_REGION + ".amazonaws.com/" + conf.BUCKET_NAME + "/" + rest_id + "/" + file.Filename
	return domain, gin.H{"status": http.StatusOK, "message": "OK"}
}
