package controllers

import (
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

func RestaurantUploadFile(c *gin.Context) (string, gin.H) {
	file, err := c.FormFile("file")
	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": "Could not retrive the file from the request!"}
	}
	rest_id := c.PostForm("rest_id")

	filePath := "/www/menu-icons/" + rest_id + "/" + file.Filename
	err = c.SaveUploadedFile(file, filePath)

	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": err.Error()}
	}
	domain, header := UploadToS3(filePath)
	rid, _ := strconv.ParseInt(rest_id, 10, 64)
	db := CreateConnection()
	results, err := db.Query("UPDATE restaurants SET logo = (?) WHERE rest_id = (?)", domain, rid)
	CloseConnection(db)
	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": "Insertion Error! Create Product", "error": err.Error()}
	}
	return domain, gin.H{"status": http.StatusOK, "message": "OK", "results": results, "header": header}
}

func ProductUploadFile(c *gin.Context) (string, gin.H) {
	file, err := c.FormFile("file")
	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": "Could not retrive the file from the request!", "error": err.Error()}
	}
	rest_id := c.PostForm("rest_id")
	prod_id := c.PostForm("prod_id")

	filePath := "/www/menu-icons/" + rest_id + "/" + prod_id + "/" + file.Filename
	err = c.SaveUploadedFile(file, filePath)

	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": err.Error()}
	}

	domain, header := UploadToS3(filePath)
	rid, _ := strconv.ParseInt(rest_id, 10, 64)
	pid, _ := strconv.ParseInt(prod_id, 10, 64)
	db := CreateConnection()
	results, err := db.Query("UPDATE products SET prod_image = (?) WHERE rest_id = (?) AND prod_id = (?)", domain, rid, pid)
	CloseConnection(db)
	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": "Insertion Error! Product Upload File", "error": err.Error()}
	}
	return domain, gin.H{"status": http.StatusOK, "message": "OK", "results": results, "header": header}
}

func CategoryUploadFile(c *gin.Context) (string, gin.H) {
	file, err := c.FormFile("file")
	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": "Could not retrive the file from the request!", "error": err.Error()}
	}
	rest_id := c.PostForm("rest_id")
	cat_id := c.PostForm("cat_id")

	filePath := "/www/menu-icons/" + rest_id + "/" + cat_id + "/" + file.Filename
	err = c.SaveUploadedFile(file, filePath)

	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": err.Error()}
	}

	domain, header := UploadToS3(filePath)
	rid, _ := strconv.ParseInt(rest_id, 10, 64)
	cid, _ := strconv.ParseInt(cat_id, 10, 64)
	db := CreateConnection()
	results, err := db.Query("UPDATE categories SET cat_image = (?) WHERE rest_id = (?) AND cat_id = (?)", domain, rid, cid)
	CloseConnection(db)
	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": "Insertion Error! Category Upload File", "error": err.Error()}
	}
	return domain, gin.H{"status": http.StatusOK, "message": "OK", "results": results, "header": header}
}

func UploadToS3(filePath string) (string, gin.H) {
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
		Key:    aws.String(f.Name()),
		ACL:    aws.String("private"),
		Body:   f,
	})
	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": err.Error()}
	}
	domain := "https://s3." + conf.BUCKET_REGION + ".amazonaws.com/" + conf.BUCKET_NAME + f.Name()
	return domain, gin.H{"status": http.StatusOK, "message": "OK"}
}
