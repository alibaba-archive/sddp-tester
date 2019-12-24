package lib

import (
	"../common"
	"fmt"
	"os"
)

func Abnormal() {
	var std *os.File

	os.Stdout = std
	os.Stderr = std

	if Config["ODPS"].(bool) {
		//配置失当-ODPS敏感项目未设置保护
		Odps("set ProjectProtection=FALSE;")

		//配置失当-ODPS敏感项目未设置标签安全
		Odps("Set LabelSecurity=true;")
		Odps("SET LABEL 0 TO TABLE " + Config["STable"].(string) + ";")
		Odps("select * from " + Config["STable"].(string) + ";")

		//异常地址、ip
		if Config["PROXY"].(string) != "" {
			os.Setenv("HTTP_PROXY", Config["PROXY"].(string))
			os.Setenv("HTTPS_PROXY", Config["PROXY"].(string))
			Odps("select * from " + Config["STable"].(string) + ";")
			os.Unsetenv("HTTP_PROXY")
			os.Unsetenv("HTTPS_PROXY")
		}

		//下载量偏高异常
		for i := 0; i <= 40; i++ {
			Odps("select * from " + Config["STable"].(string) + ";")
		}

	}

	if Config["OSS"].(bool) {
		//配置失当-OSS敏感Bucket被设置为公开
		OSSPutAcl(3)

		//多次尝试访问未成功(403、404)
		OSSBruteForce()

		//初次下载敏感数据
		OSSGetFile(Config["SObject"].(string), "test.jpg")
		os.Remove("test.jpg")

		//异常地址、ip
		if Config["PROXY"].(string) != "" {
			os.Setenv("HTTP_PROXY", Config["PROXY"].(string))
			os.Setenv("HTTPS_PROXY", Config["PROXY"].(string))
			OSSGetFile(Config["SObject"].(string), "test.jpg")
			os.Unsetenv("HTTP_PROXY")
			os.Unsetenv("HTTPS_PROXY")
			os.Remove("test.jpg")
		}

		//下载量偏高异常
		Objects, err := OSSGetBucketObject()
		if err == nil {
			for _, object := range Objects {
				fmt.Print(object)
				OSSGetFile(object, "test")
				os.Remove("test")
			}
		}

	}
	if Config["RDS"].(bool) {
		//异常地址、ip
		if Config["PROXY"].(string) != "" {
			os.Setenv("HTTP_PROXY", Config["PROXY"].(string))
			os.Setenv("HTTPS_PROXY", Config["PROXY"].(string))
			RdsQuery("select * from " + Config["STable"].(string) + ";")
			os.Unsetenv("HTTP_PROXY")
			os.Unsetenv("HTTPS_PROXY")
		}
		//下载量偏高异常
		for i := 0; i <= 40; i++ {
			RdsQuery("select * from " + Config["STable"].(string) + ";")
		}

	}
}

func OSSBruteForce() {
	AccessKeySecret := Config["AccessKeySecret"]
	for i := 0; i <= 100; i++ {
		Config["AccessKeySecret"] = common.GetRandomString(30)
		OSSGetFile(Config["SObject"].(string), "test.jpg")
	}

	Config["AccessKeySecret"] = AccessKeySecret

	for i := 0; i <= 100; i++ {
		OSSGetFile("404.JPG", "test.jpg")
	}
}
