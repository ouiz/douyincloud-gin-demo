package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// 定义一个结构体，用于存储配置信息
type Config struct {
	// DBHost string `json:"db_host"` // 数据库的主机地址
	// DBPort int    `json:"db_port"` // 数据库的端口号
	// DBUser string `json:"db_user"` // 数据库的用户名
	// DBPass string `json:"db_pass"` // 数据库的密码
	// APIKey string `json:"api_key"` // API的密钥
	AppId   string `json:"appId"`
	Secret  string `json:"secret"`
	IsLocal bool   `json:"isLocal"`

	OsAK string `json:"osAK"`
	OsSK string `json:"osSK"`
	//  //填写从抖音云--对象存储--配置获取的AK和SK
	//  osAK  = "xxxxxxxxxx"
	//  osSK  = "xxxxxx"
	//  // 如果部署在抖音云服务中，建议替换成内网域名
	//  osEndpoint   = "tos-cn-beijing.volces.com"
	//  osRegion     = "cn-beijing"
	//  // 填写从抖音云--对象存储获取的桶的ID
	//  osBucketName = "xxxxxxxxxx"
	//  httpClient = &http.Client{}
}
type ConfigCommon struct {
	OsEndpoint   string `json:"osEndpoint"`
	OsRegion     string `json:"osRegion"`
	OsBucketName string `json:"osBucketName"`
	//  // 如果部署在抖音云服务中，建议替换成内网域名
	//  osEndpoint   = "tos-cn-beijing.volces.com"
	//  osRegion     = "cn-beijing"
	//  // 填写从抖音云--对象存储获取的桶的ID
	//  osBucketName = "xxxxxxxxxx"
}

var Cfg *Config
var CfgCm = ConfigCommon{
	// OsEndpoint:   "tt278e5133b8bd3b8b01-env-nkjyi3kslh.tos-cn-beijing.volces.com",
	OsEndpoint:   "tos-cn-beijing.volces.com",
	OsRegion:     "cn-beijing",
	OsBucketName: "tt278e5133b8bd3b8b01-env-nkjyi3kslh",
}

// 定义一个全局变量，存储配置文件的名称
const configFileName = "config.json"

// 定义一个函数，用于读取配置文件，返回配置信息
func readConfig() (*Config, error) {
	var config Config
	// 打开配置文件
	file, err := os.Open(configFileName)
	if err != nil {
		// 如果打开失败，返回错误
		return nil, err
	}
	// 延迟关闭配置文件
	defer file.Close()
	// 创建一个JSON解码器
	decoder := json.NewDecoder(file)
	// 将配置文件的内容解码为配置信息
	err = decoder.Decode(&config)
	if err != nil {
		// 如果解码失败，返回错误
		return nil, err
	}
	// 返回配置信息
	return &config, nil
}

// 定义一个函数，用于检查配置文件是否存在，返回布尔值
func checkConfigFile() bool {
	// 使用os包中的Stat函数，获取配置文件的信息
	_, err := os.Stat(configFileName)
	if err != nil {
		// 如果出错，判断是否是因为配置文件不存在
		if os.IsNotExist(err) {
			// 如果配置文件不存在，返回false
			return false
		}
		// 如果是其他原因，打印错误并退出程序
		log.Fatal(err)
	}
	// 如果没有出错，表示配置文件存在，返回true
	return true
}

var cldCfg = Config{
	AppId:  os.Getenv("APPID"),
	Secret: os.Getenv("SECRET"),
	OsAK:   os.Getenv("OS_AK"),
	OsSK:   os.Getenv("OS_SK"),
}

func InitCfg() {
	var err error
	if checkConfigFile() {
		// 如果配置文件存在，读取配置文件，获取配置信息
		Cfg, err = readConfig()
		if err != nil {
			// log.Fatal(err)
			fmt.Println("readCfg err", err)
			Cfg = &cldCfg
		}
		// 打印配置信息
	} else {
		Cfg = &cldCfg
	}
	fmt.Println("Config init after:", Cfg)
}
