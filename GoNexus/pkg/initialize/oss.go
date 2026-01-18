package initialize

import (
	"fmt"
	"go-nexus/pkg/global"
	"go-nexus/pkg/utils/upload"
)

func InitOSS() {
	fmt.Println(global.Config.Upload.Oss.AccessKey)
	fmt.Println(global.Config.Upload.Oss.SecretKey)
	// 根据配置文件，只在启动时创建一次实例
	switch global.Config.Upload.Type {
	case "oss":
		// 赋值给全局变量
		global.OSS = &upload.AliyunOSS{}
	case "local":
		// global.OSS = &upload.LocalOSS{} // 如果有本地上传实现
	default:
		global.OSS = &upload.AliyunOSS{} // 默认
	}
}
