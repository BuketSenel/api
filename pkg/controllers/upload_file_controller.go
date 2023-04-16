package controllers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

const (
	AWS_S3_REGION = ""
	AWS_S3_BUCKET = ""
)

func UploadFile(c *gin.Context) (bool, gin.H) {
	session, err := session.NewSession(&aws.Config{Region: aws.String(AWS_S3_REGION)})
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Error creating session!"}
	}

	// Upload Files
	err = addFilesToS3(session, "test.png")
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Error uploading file!"}
	}

	return true, gin.H{"status": http.StatusOK, "message": "File uploaded successfully!"}
}

func addFilesToS3(session *session.Session, uploadFileDir string) error {

	upFile, err := os.Open(uploadFileDir)
	if err != nil {
		return err
	}
	defer upFile.Close()

	upFileInfo, _ := upFile.Stat()
	var fileSize int64 = upFileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	upFile.Read(fileBuffer)

	_, err = s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(AWS_S3_BUCKET),
		Key:                  aws.String(uploadFileDir),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(fileBuffer),
		ContentLength:        aws.Int64(fileSize),
		ContentType:          aws.String(http.DetectContentType(fileBuffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}

// Bu da chatgpt'nin yazdığı kod ben çok beğenmedim
func main() {
	// Initialize the Gin-Gonic router
	r := gin.Default()

	// Define a route that uploads a file to S3
	r.POST("/upload", func(c *gin.Context) {
		// Get the file from the request
		file, err := c.FormFile("file")
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// Open the file
		f, err := file.Open()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer f.Close()

		// Read the file contents
		fileContents, err := ioutil.ReadAll(f)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Upload the file to S3
		sess := session.Must(session.NewSession(&aws.Config{
			Region: aws.String("us-west-2"),
		}))
		svc := s3.New(sess)
		_, err = svc.PutObject(&s3.PutObjectInput{
			Bucket: aws.String("your-bucket-name"),
			Key:    aws.String(file.Filename),
			Body:   bytes.NewReader(fileContents),
		})
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Send the response
		c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
	})

	// Run the server
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
