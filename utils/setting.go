package utils

import (
	"fmt"
	"gopkg.in/ini.v1"
)

var (
	AppMode     string
	HttpPort    string
	JwKey       string
	Zone        int
	Db          string
	DbHost      string
	DbPort      string
	DbUser      string
	DbPassWord  string
	DbName      string
	AccessKey   string
	SecretKey   string
	Bucket      string
	QiniuServer string
)

func init() {
	file, err := ini.Load("config/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径", err)
	}
	LoadServer(file)
	LoadData(file)
	LoadQiniu(file)
}
func LoadServer(file *ini.File) {
	AppMode = file.Section("server").Key("AppMode").MustString("debug")
	//获取不到key的值 自动返回 后面的默认值
	HttpPort = file.Section("server").Key("HttpPort").MustString("3.1.135.132:3030") //链接ip端口
	JwKey = file.Section("server").Key("JwKey").MustString("asd45645ad4cergd1")      //链接ip端口

}
func LoadData(file *ini.File) {
	Db = file.Section("database").Key("Db").MustString("mysql")
	DbHost = file.Section("database").Key("DbHost").MustString("") //数据库ip

	DbPort = file.Section("database").Key("DbPort").MustString("3307")

	DbUser = file.Section("database").Key("DbUser").MustString("root")

	DbPassWord = file.Section("database").Key("DbPassWord").MustString("lsqA777..")
	DbName = file.Section("database").Key("DbName").MustString("bbs-go")

}
func LoadQiniu(file *ini.File) {
	Zone = file.Section("qiniu").Key("Zone").MustInt(1)
	AccessKey = file.Section("qiniu").Key("AccessKey").String()
	SecretKey = file.Section("qiniu").Key("SecretKey").String()
	Bucket = file.Section("qiniu").Key("Bucket").String()
	QiniuServer = file.Section("qiniu").Key("QiniuServer").String()

}
