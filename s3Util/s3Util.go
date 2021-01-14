package s3Util

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gyf841010/pz-infra-new/log"

	"github.com/astaxie/beego"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type UploadFileContentInput struct {
	//上传到aws的文件路径 例: a/b/c/1.txt
	UploadPath string
	//文件内容
	FileContent []byte
	//是否转换为HTTP
	ConvertHTTP bool
}

func UploadFileContent(input *UploadFileContentInput) (string, error) {
	if input.FileContent == nil {
		return "", errors.New("file content is emtpy !")
	}
	awsAccessKeyId := beego.AppConfig.String("awsAccessKeyId")
	awsSecretAccessKey := beego.AppConfig.String("awsSecretAccessKey")
	if awsAccessKeyId == "" || awsSecretAccessKey == "" {
		return "", log.Error("config of awsAccessKeyId or awsSecretAccessKey is empty")
	}
	s3Bucket := beego.AppConfig.String("awsS3Bucket")
	if s3Bucket == "" {
		return "", log.Error("aws bucket is empty, please check your configuration ...")
	}

	sess := session.New(&aws.Config{
		Region:      aws.String(beego.AppConfig.DefaultString("awsS3Region", "us-east-1")),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyId, awsSecretAccessKey, ""),
	})

	mgr := s3manager.NewUploader(sess)

	contentType := http.DetectContentType(input.FileContent)

	resp, err := mgr.Upload(&s3manager.UploadInput{
		Bucket:               aws.String(s3Bucket),
		Key:                  aws.String(input.UploadPath),
		Body:                 bytes.NewReader(input.FileContent),
		ServerSideEncryption: aws.String("AES256"),
		ContentType:          aws.String(contentType),
	})

	if err != nil {
		log.Info("upload ", input.UploadPath, " failed", err)
		return "", nil
	}

	if input.ConvertHTTP {
		convertedUrl := strings.Replace(resp.Location, "https://", "http://", 1)
		return convertedUrl, nil
	}
	return resp.Location, nil
}

type UploadFileInput struct {
	//上传到aws的文件路径 例: a/b/c/1.txt
	UploadPath string
	//文件的本地路径
	LocalFilePath string
}

func UploadFile(input *UploadFileInput) (string, error) {
	fileContent, err := ioutil.ReadFile(input.LocalFilePath)
	if err != nil {
		return "", err
	}
	return UploadFileContent(&UploadFileContentInput{
		FileContent: fileContent,
		UploadPath:  input.UploadPath,
	})
}

func DownloadFile(url, localPath string) error {
	out, err := os.Create(localPath)
	if err != nil {
		log.Error("failed to create file ", localPath, err)
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		log.Error("failed to download from ", url, err)
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Error("failed to copy file content", err)
		return err
	}
	return nil
}

func getS3Bucket() string {
	return beego.AppConfig.String("awsS3Bucket")
}
