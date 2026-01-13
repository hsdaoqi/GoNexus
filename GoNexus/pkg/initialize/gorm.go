package initialize

import (
	"fmt"
	"go-nexus/pkg/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func InitMySQL() {
	m := global.Config.Mysql

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.UserName, m.Password, m.Host, m.Port, m.DbName)

	// 1. 根据配置决定日志级别
	// 开发环境(debug)：打印所有 SQL，方便调试
	// 生产环境(release)：只打印 Error，避免日志刷屏影响性能
	var logMode logger.LogLevel
	if global.Config.Server.Mode == "debug" {
		logMode = logger.Info
	} else {
		logMode = logger.Error
	}
	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{

		// 用法：设置日志级别 (Silent, Error, Warn, Info)。
		Logger: logger.Default.LogMode(logMode),
	})

	if err != nil {
		panic(fmt.Errorf("MySQL连接失败：%s", err))
	}
	fmt.Println("MySQL数据库连接成功")

	// 2. 配置连接池 (Connection Pooling) —— 复试核心考点！
	sqlDB, _ := global.DB.DB() // 获取底层的 *sql.DB 对象

	// SetMaxIdleConns: 设置空闲连接池中连接的最大数量
	// 解释：哪怕没人用，我也保持 10 个连接连着数据库，这样下一个请求来的时候就不用重新握手，直接用。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns: 设置打开数据库连接的最大数量
	// 解释：同时最多允许 100 个人连数据库。如果超过 100，第 101 个人就要排队。
	// 防止流量激增把数据库压垮。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime: 设置了连接可复用的最大时间
	// 解释：一个连接用了一个小时，强制关掉重新连一下，防止 MySQL 服务器那边超时断开导致报错。
	sqlDB.SetConnMaxLifetime(time.Hour)
}
