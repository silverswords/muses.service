package main

import (
	"fmt"
	"os"

	"github.com/Unknwon/goconfig"
	"muses.service/apis"
)

var cfg *goconfig.ConfigFile

// InitConf -
func InitConf() {
	config, err := goconfig.LoadConfigFile("conf/database.conf") //加载配置文件
	if err != nil {
		fmt.Println("get config file error")
		os.Exit(-1)
	}

	cfg = config
}

func main() {
	InitConf()
	r := apis.InitRouter()

	port, _ := cfg.GetValue("gin", "port")
	ginServer := fmt.Sprintf(":%v", port)

	r.Run(ginServer)
}
