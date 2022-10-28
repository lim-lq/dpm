package app

import (
	"flag"

	"github.com/lim-lq/dpm/core"
	"github.com/lim-lq/dpm/core/config"
	"github.com/lim-lq/dpm/core/log"
)

func Initialize() {
	// 解析参数
	var configFile string
	flag.StringVar(&configFile, "c", "config.yaml", "config file path")
	flag.Parse()

	// 初始化配置文件
	config.InitConfig(configFile)

	// 初始化日志
	log.InitLogger()

	// 初始化redis
	core.InitRedis()

	// 初始化mongodb
	core.InitMongo()
}

func Run() {
	Initialize()
	log.Logger.Info("Begin run dpm")
	// 启动服务
	log.Logger.Fatal(StartServer())
}
