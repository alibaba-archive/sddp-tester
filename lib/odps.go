package lib

import (
	"../common"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func OdpsCreate(command string) (*http.Response, error) {
	url := Config["ODPSEndPoint"].(string) + "/projects/" + Config["ODPSProject"].(string) + "/instances?curr_project=" + Config["ODPSProject"].(string)
	data := fmt.Sprintf("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n<Instance>\n  <Job>\n    <Priority>9</Priority>\n    <Tasks>\n      <SQL>\n        <Name>AnonymousSQLTask</Name>\n        <Config>\n          <Property>\n            <Name>settings</Name>\n            <Value>{}</Value>\n          </Property>\n          <Property>\n            <Name>uuid</Name>\n            <Value>%s</Value>\n          </Property>\n        </Config>\n        <Query><![CDATA[%s;]]></Query>\n      </SQL>\n    </Tasks>\n    <DAG>\n      <RunMode>Sequence</RunMode>\n    </DAG>\n  </Job>\n</Instance>\n", common.TimeUUID(), command)
	timestr := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	lines := []string{"POST", "", "application/xml", timestr, "/projects/" + Config["ODPSProject"].(string) + "/instances?curr_project=" + Config["ODPSProject"].(string)}
	hash := hmac.New(sha1.New, []byte(Config["AccessKeySecret"].(string)))
	hash.Write([]byte([]byte(strings.Join(lines, "\n"))))
	authorization := "ODPS " + Config["AccessKeyId"].(string) + ":" + base64.StdEncoding.EncodeToString(hash.Sum(nil))
	client := &http.Client{}
	reqest, _ := http.NewRequest("POST", url, strings.NewReader(data))
	reqest.Header.Add("Accept-Encoding", "identity")
	reqest.Header.Add("User-Agent", "pyodps/0.8.1 CPython/3.7.0 Darwin/18.2.0")
	reqest.Header.Add("Content-Type", "application/xml")
	reqest.Header.Add("Authorization", authorization)
	reqest.Header.Add("Date", timestr)
	return client.Do(reqest)
}

func Odps(sqlcommand string) error {
	response, _ := OdpsCreate(sqlcommand)
	if response.StatusCode == 201 {
		fmt.Println("Create success")
		return nil
	} else {
		return fmt.Errorf("Exec failed,Check your project")
	}
}

func OdpsGetTableName() {
	Config["TableName"] = common.GetRandomString(10)
}
