package configs

type Server struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type Mysql struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	UserName string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"dbname"`
}

type Jwt struct {
	Secret string `mapstructure:"secret"`
	Expire int64  `mapstructure:"expire"`
}

type Upload struct {
	Type string `mapstructure:"type"`
	Path string `mapstructure:"path"`
	Oss  Oss    `mapstructure:"oss"`
}

type Oss struct {
	EndPoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	Domain    string `mapstructure:"domain"`
}

// 总配置结构体
type Config struct {
	Server Server `mapstructure:"server"`
	Mysql  Mysql  `mapstructure:"mysql"`
	Jwt    Jwt    `mapstructure:"jwt"`
	Upload Upload `mapstructure:"upload"`
}
