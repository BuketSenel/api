package controllers

import (
	"net/http"
	"os"
	"strconv"

	"github.com/SelfServiceCo/api/pkg/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) (string, gin.H) {
	file, err := c.FormFile("file")
	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": "Error uploading file!"}
	}
	products := models.Product{}

	if err := c.BindJSON(&products); err != nil {
		c.AbortWithError(401, err)
	}
	filePath := "/www/uploads/" + strconv.FormatInt(int64(products.ResID), 10) + "/" + file.Filename
	c.SaveUploadedFile(file, filePath)

	f, err := os.Open(filePath)
	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": "Error uploading file!"}
	}
	session, err := session.NewSession(&aws.Config{Region: aws.String(conf.BUCKET_REGION), Credentials: credentials.NewStaticCredentials(conf.BUCKET_KEY, conf.BUCKET_SECRET, "")})
	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": "Error creating session!"}
	}
	_, err = s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(conf.BUCKET_NAME),
		Key:    aws.String(strconv.FormatInt(int64(products.ResID), 10) + "/" + file.Filename),
		ACL:    aws.String("private"),
		Body:   f,
	})
	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": "Error uploading file!"}
	}
	return filePath, gin.H{"status": http.StatusOK, "message": "OK"}
}
