package lib

import (
	"../common"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"strings"
)

func OSSPutDir(ObjectDirName, PathName string) error {
	client, err := oss.New(Config["OSSEndPoint"].(string), Config["AccessKeyId"].(string), Config["AccessKeySecret"].(string))
	if err != nil {
		return fmt.Errorf("Open Connection failed:", err.Error())
	}

	bucket, err := client.Bucket(Config["OSSBucketName"].(string))
	if err != nil {
		return fmt.Errorf("Open Bucket failed:", err.Error())
	}
	var FileNames []string
	FileNames, err = common.GetAllFile(PathName, FileNames)
	if err != nil {
		return fmt.Errorf("Read Dir failed:", err.Error())
	}
	for _, FileName := range FileNames {
		ObjectFileName := strings.TrimRight(ObjectDirName, "/") + strings.Replace(FileName, PathName, "", 1)
		err = bucket.PutObjectFromFile(ObjectFileName, FileName)
		Config["Temp"] = Config["Temp"].(string) + "PutFile:" + ObjectFileName + "\n"
		if err != nil {
			fmt.Println("Upload Failed:", err.Error())
			continue
		}
	}
	fmt.Println("Put Success")
	return nil

}

func OSSPutFile(ObjectName string, FileName string) error {
	client, err := oss.New(Config["OSSEndPoint"].(string), Config["AccessKeyId"].(string), Config["AccessKeySecret"].(string))
	if err != nil {
		return fmt.Errorf("Open Connection failed:", err.Error())
	}

	bucket, err := client.Bucket(Config["OSSBucketName"].(string))
	if err != nil {
		return fmt.Errorf("Open Bucket failed:", err.Error())
	}

	err = bucket.PutObjectFromFile(ObjectName, FileName)
	if err != nil {
		return fmt.Errorf("Upload Failed:", err.Error())
	}
	fmt.Println("Put Success")
	return nil
}

func OSSGetFile(ObjectName string, FileName string) error {
	client, err := oss.New(Config["OSSEndPoint"].(string), Config["AccessKeyId"].(string), Config["AccessKeySecret"].(string))
	if err != nil {
		return fmt.Errorf("Open Connection failed:", err.Error())
	}

	bucket, err := client.Bucket(Config["OSSBucketName"].(string))
	if err != nil {
		return fmt.Errorf("Open Bucket failed:", err.Error())
	}

	err = bucket.GetObjectToFile(ObjectName, FileName)
	if err != nil {
		return fmt.Errorf("Get Failed:", err.Error())
	}
	fmt.Println("Get Success")
	return nil

}

func OSSDelete(ObjectName string) error {
	client, err := oss.New(Config["OSSEndPoint"].(string), Config["AccessKeyId"].(string), Config["AccessKeySecret"].(string))
	if err != nil {
		return fmt.Errorf("Open Connection failed:", err.Error())
	}
	bucket, err := client.Bucket(Config["OSSBucketName"].(string))
	if err != nil {
		return fmt.Errorf("Open Bucket failed:", err.Error())
	}
	err = bucket.DeleteObject(ObjectName)
	if err != nil {
		return fmt.Errorf("Delete Failed:", err.Error())
	}
	fmt.Println("Delete Success")
	return nil
}

func OSSPutAcl(AclInt int) error {
	client, err := oss.New(Config["OSSEndPoint"].(string), Config["AccessKeyId"].(string), Config["AccessKeySecret"].(string))
	if err != nil {
		return fmt.Errorf("Open Connection failed:", err.Error())
	}
	var Acl oss.ACLType
	switch AclInt {
	case 1:
		{
			Acl = oss.ACLPrivate
		}
	case 2:
		{
			Acl = oss.ACLPublicRead
		}
	case 3:
		{
			Acl = oss.ACLPublicReadWrite
		}
	case 4:
		{
			Acl = oss.ACLDefault
		}
	}

	err = client.SetBucketACL(Config["OSSBucketName"].(string), Acl)
	if err != nil {
		return fmt.Errorf("Put Acl failed:", err.Error())
	}
	fmt.Println("Acl Set Success")
	return nil
}

func OSSGetDir() error {
	client, err := oss.New(Config["OSSEndPoint"].(string), Config["AccessKeyId"].(string), Config["AccessKeySecret"].(string))
	if err != nil {
		return fmt.Errorf("Open Connection failed:", err.Error())
	}
	bucket, err := client.Bucket(Config["OSSBucketName"].(string))
	if err != nil {
		return fmt.Errorf("Open Bucket failed:", err.Error())
	}
	for {
		dir := common.GetRandomString(10)
		lsRes, err := bucket.ListObjects(oss.Prefix(dir))
		if err != nil {
			return fmt.Errorf("File Exist Check failed:", err.Error())
		} else {
			if len(lsRes.Objects) == 0 {
				Config["OSSDir"] = dir
				return nil
			}
		}
	}
}

func OSSGetBucketObject() ([]string, error) {
	var Objects []string
	client, err := oss.New(Config["OSSEndPoint"].(string), Config["AccessKeyId"].(string), Config["AccessKeySecret"].(string))
	if err != nil {
		return nil, fmt.Errorf("Open Connection failed:", err.Error())
	}
	bucket, err := client.Bucket(Config["OSSBucketName"].(string))
	if err != nil {
		return nil, fmt.Errorf("Open Bucket failed:", err.Error())
	}

	lsRes, err := bucket.ListObjects(oss.MaxKeys(1000))
	if err != nil {
		return nil, fmt.Errorf("Get Objects failed:", err)
	}
	for _, object := range lsRes.Objects {
		Objects = append(Objects, object.Key)
	}
	return Objects, nil
}
