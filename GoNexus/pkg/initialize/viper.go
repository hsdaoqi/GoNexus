package initialize

import (
	"fmt"
	"github.com/spf13/viper"
	"go-nexus/pkg/global"
	"strings"
)

func InitConfig() {
	v := viper.New()
	v.SetConfigFile("configs/config.yaml") //指定配置文件路径
	// 1. 设置环境变量前缀 (建议用项目名大写，防止冲突)
	// 比如 config.yaml 里的 mysql.host，环境变量就是 GONEXUS_MYSQL_HOST
	v.SetEnvPrefix("GONEXUS")
	// 2. 将配置文件的点号 "." 替换为下划线 "_"
	// 因为环境变量通常不支持点号
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// 3. 开启自动读取
	// Viper 会自动查找匹配的环境变量并覆盖 config.yaml 里的值
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败：%s", err))
	}

	// 把读取到的配置映射到全局变量 global.Config 中
	if err := v.Unmarshal(&global.Config); err != nil {
		panic(fmt.Errorf("配置解析失败：%s", err))
	}
}
