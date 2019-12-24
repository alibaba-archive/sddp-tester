package main

import (
	"./common"
	"./lib"
	"bytes"
	"fmt"
	"github.com/urfave/cli"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	rdstype         string
	host            string
	port            string
	database        string
	endpoint        string
	bucket          string
	project         string
	product         string
	sensitiveobject string
	proxy           string
	table           string
)

func main() {
	oApp := cli.NewApp()
	oApp.Name = "SDDP-tester"
	oApp.Usage = "To test the SDDP"
	oApp.Version = "1.0.0"
	oApp.UsageText = `main command [command options] [arguments...]

Example:
	./sddp-tester scan -p rds -t mysql  -d test -P 3306 -H xxx.mysql.rds.aliyuncs.com 
	./sddp-tester scan -p oss  --endpoint http://oss-cn-beijing.aliyuncs.com --bucket sddp_test 
	./sddp-tester scan -p odps --project sddp_test  
	./sddp-tester anomaly -p rds -t mysql  -d test --table test -P 3306 -H xxx.mysql.rds.aliyuncs.com --proxy http://1.1.1.1:8888 
	./sddp-tester anomaly -p oss  --endpoint http://oss-cn-beijing.aliyuncs.com --bucket test --sensitiveobject test/1.jpg --proxy http://1.1.1.1:8888  
	./sddp-tester anomaly -p odps --project test  --table test --proxy http://1.1.1.1:8888 
	./sddp-tester clean
`

	oApp.Commands = []cli.Command{
		{
			Name:  "clean",
			Usage: "Clean up test files",
			Action: func(c *cli.Context) {
				Clean()
			},
		},
		{
			Name: "anomaly",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "product,p",
					Value:       "",
					Usage:       "Required: <odps or rds or oss>",
					Destination: &product,
				},
				cli.StringFlag{
					Name:        "rdstype,t",
					Value:       "",
					Usage:       "Required (RDS only):  <mysql or mssql>",
					Destination: &rdstype,
				},
				cli.StringFlag{
					Name:        "host,H",
					Value:       "",
					Usage:       "Required (RDS only):  <RDS hostname [Example: xxx.mysql.rds.aliyuncs.com]>",
					Destination: &host,
				},
				cli.StringFlag{
					Name:        "port,P",
					Value:       "3306",
					Usage:       "Required (RDS only):  <port number [default:3306]>",
					Destination: &port,
				},
				cli.StringFlag{
					Name:        "database,d",
					Value:       "",
					Usage:       "Required (RDS only): <RDS database name [Example: database_test]>",
					Destination: &database,
				},
				cli.StringFlag{
					Name:        "endpoint",
					Value:       "",
					Usage:       "Required (OSS only): <OSS endpoint [Example: http://oss-cn-beijing.aliyuncs.com]> ",
					Destination: &endpoint,
				},
				cli.StringFlag{
					Name:        "bucket",
					Value:       "",
					Usage:       "Required (OSS only): <OSS bucket name [Example: MyBucket]>",
					Destination: &bucket,
				},
				cli.StringFlag{
					Name:        "project",
					Value:       "",
					Usage:       "Required (MaxCompute only): <MaxCompute project name[Example:MyProject]>",
					Destination: &project,
				},
				cli.StringFlag{
					Name:        "sensitiveobject,s",
					Value:       "",
					Usage:       "Optional: <Sensitive OSS object name [Example: id_card.jpeg]>",
					Destination: &sensitiveobject,
				},
				cli.StringFlag{
					Name:        "proxy",
					Value:       "",
					Usage:       "Optional: <anomalous proxy access [Example: 10.10.10.10:8080]>",
					Destination: &proxy,
				},
				cli.StringFlag{
					Name:        "table",
					Value:       "",
					Usage:       "Optional: <Sensitive RDS/ODPS table name [Example: Mytable]>",
					Destination: &table,
				},
			},
			Usage: "Generates anomalous events and Sensitive Data Discovery and Protection alarms",
			Action: func(c *cli.Context) {
				if product == "" {
					cli.ShowCommandHelp(c, "anomaly")
					return
				}
				err := CheckInput()
				if err != nil {
					fmt.Println(err)
					return
				}

				if proxy == "" {
					fmt.Fprintf(os.Stderr, "Testing in progress… Generate anomaly alerts with http proxy \"—proxy\" \n")
					lib.Config["PROXY"] = ""
				} else {
					lib.Config["PROXY"] = proxy
				}

				if lib.Config["OSS"].(bool) {
					if sensitiveobject == "" {
						fmt.Fprintf(os.Stderr, "Testing in progress… Generate anomaly alerts with sensitive objects \"—sensitiveobjects\" \n")
						lib.Config["SObject"] = ""
					} else {
						lib.Config["SObject"] = sensitiveobject
					}

				}

				if lib.Config["RDS"].(bool) {
					if table == "" {
						fmt.Fprintf(os.Stderr, "Testing in progress… Generate anomaly alerts with sensitive table \"—table\" \n")
						lib.Config["STable"] = "STable"
					} else {
						lib.Config["STable"] = table
					}
				}

				if lib.Config["ODPS"].(bool) {
					if table == "" {
						fmt.Fprintf(os.Stderr, "Testing in progress… Generate anomaly alerts with sensitive table \"—table\" \n")
						lib.Config["STable"] = "STable"
					} else {
						lib.Config["STable"] = table
					}
				}
				lib.Abnormal()
			},
		},
		{
			Name: "scan",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "product,p",
					Value:       "",
					Usage:       "Required: <odps or rds or oss>",
					Destination: &product,
				},
				cli.StringFlag{
					Name:        "rdstype,t",
					Value:       "",
					Usage:       "Required (RDS only):  <mysql or mssql>",
					Destination: &rdstype,
				},
				cli.StringFlag{
					Name:        "host,H",
					Value:       "",
					Usage:       "Required (RDS only):  <RDS hostname [Example: xxx.mysql.rds.aliyuncs.com]>",
					Destination: &host,
				},
				cli.StringFlag{
					Name:        "port,P",
					Value:       "3306",
					Usage:       "Required (RDS only):  <port number [default:3306]>",
					Destination: &port,
				},
				cli.StringFlag{
					Name:        "database,d",
					Value:       "",
					Usage:       "Required (RDS only): <RDS database name [Example: database_test]>",
					Destination: &database,
				},
				cli.StringFlag{
					Name:        "endpoint",
					Value:       "",
					Usage:       "Required (OSS only): <OSS endpoint [Example: %RANAN%]> ",
					Destination: &endpoint,
				},
				cli.StringFlag{
					Name:        "bucket",
					Value:       "",
					Usage:       "Required (OSS only): <OSS bucket name [Example: MyBucket]>",
					Destination: &bucket,
				},
				cli.StringFlag{
					Name:        "project",
					Value:       "",
					Usage:       "Required (MaxCompute only): <MaxCompute project name[Example:MyProject]>",
					Destination: &project,
				},
			},
			Usage: "Sensitive data scan",
			Action: func(c *cli.Context) {
				if product == "" {
					cli.ShowCommandHelp(c, "scan")
					return
				}
				err := CheckInput()
				if err != nil {
					fmt.Println(err)
					return
				}
				err = Scan()
				if err != nil {
					fmt.Println(err)
					return
				}

			},
		},
	}
	err := oApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Scan() error {
	var buf []byte

	if lib.Config["RDS"].(bool) {
		err := lib.RdsGetTableName()
		if err != nil {
			return err
		}
		buf, _ = ioutil.ReadFile("./samples/test.sql")
		err = lib.RdsInsert(strings.Replace(string(buf), "{{table}}", lib.Config["TableName"].(string), 102))
		if err != nil {
			return err
		}
		fmt.Println("Sensitive content has been generated, Results will appear in url: https://yundunnext.console.aliyun.com/?p=sddp#/asset/rds in 30 minutes")
	}

	if lib.Config["ODPS"].(bool) {
		lib.OdpsGetTableName()
		buf, _ = ioutil.ReadFile("./samples/test.odps")
		err := lib.Odps(strings.Replace(string(buf), "{{table}}", common.GetRandomString(10), 1))
		if err != nil {
			return err
		}
		fmt.Println("Sensitive content has been generated, Results will appear in url: https://yundunnext.console.aliyun.com/?p=sddp#/asset/maxcompute in 30 minutes")
	}

	if lib.Config["OSS"].(bool) {
		err := lib.OSSGetDir()
		if err != nil {
			return err
		}
		//common.GetBigFile("./samples/oss_file/bigfile")
		err = lib.OSSPutDir(lib.Config["OSSDir"].(string), "./samples/oss_file")
		if err != nil {
			return err
		}
		//os.Remove("./samples/oss_file/bigfile")
		fmt.Println("Sensitive content has been generated, Results will appear in url: https://yundunnext.console.aliyun.com/?p=sddp#/asset/oss in 30 minutes")
	}

	err := common.SaveTemp(lib.Config)
	if err != nil {
		return fmt.Errorf("Temp Save failed:", err.Error())
	}

	return nil
}

func Clean() {
	dirPth := "temp"
	var FileNames []string
	FileNames, err := common.GetAllFile(dirPth, FileNames)
	if err != nil {
		fmt.Println("Temp Read failed:", err)
		return
	}
	for _, fi := range FileNames {
		if strings.HasSuffix(strings.ToUpper(fi), "INI") {
			lib.Config, err = common.ReadTemp(fi)
			if err != nil {
				continue
			}
			var reader bytes.Buffer
			reader.WriteString(lib.Config["Temp"].(string))
			for {
				line, err := reader.ReadString('\n')
				if err == io.EOF {
					break
				}
				if err != nil {
					fmt.Println(err)
				}
				line = strings.TrimSpace(line)
				arr := strings.Split(line, ":")

				key := strings.TrimSpace(arr[0])
				value := strings.TrimSpace(arr[1])
				switch key {
				case "PutFile":
					lib.OSSDelete(value)
					break
				case "RdsName":
					lib.RdsInsert("drop table " + value + ";;;")
					break
				}
			}
		}
	}

}

func CheckInput() error {
	if product == "rds" {
		lib.Config["RDS"] = true

		if rdstype == "" {
			return fmt.Errorf("RDS mode enabled, require rdstype parameter")
		} else if rdstype != "mysql" && rdstype != "mssql" {
			return fmt.Errorf("rdstype should be only mysql or mssql")
		} else {
			lib.Config["DatabaseType"] = rdstype
		}

		if lib.Config["User"] == "" {
			return fmt.Errorf("rds mode is enabled, user required in config.ini")
		}

		if lib.Config["Passwd"] == "" {
			return fmt.Errorf("rds mode is enabled, password required in config.ini")
		}

		if host == "" {
			return fmt.Errorf("rds mode is enabled, required Parameter host")
		} else {
			lib.Config["Server"] = host
		}

		if port == "" {
			return fmt.Errorf("rds mode is enabled, required Parameter port")
		} else {
			lib.Config["Port"] = port
		}

		if database == "" {
			return fmt.Errorf("rds mode is enabled,required Parameter database")
		} else {
			lib.Config["Database"] = database
		}
	} else {
		lib.Config["RDS"] = false
	}

	if product == "oss" {
		lib.Config["OSS"] = true

		if endpoint == "" {
			return fmt.Errorf("oss mode is enabled,required Parameter endpoint")
		} else {
			lib.Config["OSSEndPoint"] = endpoint
		}

		if bucket == "" {
			return fmt.Errorf("oss mode is enabled,required Parameter bucket")
		} else {
			lib.Config["OSSBucketName"] = bucket
		}

		if lib.Config["AccessKeyId"] == "" {
			return fmt.Errorf("oss mode is enabled,required enter accesskeyid in config.ini")
		}

		if lib.Config["AccessKeySecret"] == "" {
			return fmt.Errorf("oss mode is enabled,required enter accesskeysecret in config.ini")
		}
	} else {
		lib.Config["OSS"] = false
	}

	if product == "odps" {
		lib.Config["ODPS"] = true

		lib.Config["ODPSEndPoint"] = "http://service.cn.maxcompute.aliyun.com/api"

		if project == "" {
			return fmt.Errorf("odps mode is enabled,required Parameter project")
		} else {
			lib.Config["ODPSProject"] = project
		}

		if lib.Config["AccessKeyId"] == "" {
			return fmt.Errorf("odps mode is enabled, accesskeyid required in config.ini")
		}

		if lib.Config["AccessKeySecret"] == "" {
			return fmt.Errorf("odps mode is enabled, accesskeysecret required in config.ini")
		}
	} else {
		lib.Config["ODPS"] = false
	}

	if !(lib.Config["OSS"].(bool) || lib.Config["ODPS"].(bool) || lib.Config["RDS"].(bool)) {
		return fmt.Errorf("product mode is required, should be only odps or rds or oss")
	}
	return nil
}
