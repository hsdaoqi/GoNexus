package upload

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"go-nexus/pkg/global"
	"mime/multipart"
	"path"
	"time"
)

type AliyunOSS struct{}

func (a *AliyunOSS) UploadFile(file *multipart.FileHeader) (string, error) {
	cfg := global.Config.Upload.Oss

	//1.创建OSS Client
	client, err := oss.New(cfg.EndPoint, cfg.AccessKey, cfg.SecretKey)
	if err != nil {
		return "", err
	}

	//2.获取存储空间
	bucket, err := client.Bucket(cfg.Bucket)
	if err != nil {
		return "", err
	}

	//3.打开文件流
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	// 4. 生成唯一文件名 (防止重名覆盖)
	// 格式: uploads/YYYYMM/uuid.jpg
	suffix := path.Ext(file.Filename)
	fileName := fmt.Sprintf("upload/%s/%s%s",
		time.Now().Format("200601"),
		uuid.New().String(),
		suffix,
	)

	//5.上传
	err = bucket.PutObject(fileName, f)
	if err != nil {
		return "", err
	}
	//6.返回完整访问路径
	return cfg.Domain + "/" + fileName, nil
}

func (a *AliyunOSS) DeleteFile(key string) error {
	//TODO 删除文件
	return nil
}
