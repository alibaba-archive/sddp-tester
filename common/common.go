package common

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

func PrintRow(colsdata []interface{}) {
	for _, val := range colsdata {
		switch v := (*(val.(*interface{}))).(type) {
		case nil:
			fmt.Print("NULL")
		case bool:
			if v {
				fmt.Print("True")
			} else {
				fmt.Print("False")
			}
		case []byte:
			fmt.Print(string(v))
		case time.Time:
			fmt.Print(v.Format("2016-01-02 15:05:05.999"))
		default:
			fmt.Print(v)
		}
		fmt.Print("\t")
	}
	fmt.Println()
}
func GetAllFile(pathname string, filenames []string) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return filenames, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			FullDir := pathname + "/" + fi.Name()
			filenames, err = GetAllFile(FullDir, filenames)
			if err != nil {
				return filenames, err
			}
		} else {
			FullName := pathname + "/" + fi.Name()
			filenames = append(filenames, FullName)
		}
	}
	return filenames, nil
}

func GetRandomString(stringlen int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < stringlen; i++ {
		result = append(result, str[r.Intn(len([]byte(str)))])
	}
	return string(result)
}

type UUID [16]byte

var timeBase = time.Date(1582, time.October, 15, 0, 0, 0, 0, time.UTC).Unix()
var hardwareAddr []byte
var clockSeq uint32

func TimeUUID() UUID {
	return FromTime(time.Now())
}

func FromTime(atime time.Time) UUID {
	var u UUID

	utcTime := atime.In(time.UTC)
	t := uint64(utcTime.Unix()-timeBase)*10000000 + uint64(utcTime.Nanosecond()/100)
	u[0], u[1], u[2], u[3] = byte(t>>24), byte(t>>16), byte(t>>8), byte(t)
	u[4], u[5] = byte(t>>40), byte(t>>32)
	u[6], u[7] = byte(t>>56)&0x0F, byte(t>>48)

	clock := atomic.AddUint32(&clockSeq, 1)
	u[8] = byte(clock >> 8)
	u[9] = byte(clock)

	copy(u[10:], hardwareAddr)

	u[6] |= 0x10 // set version to 1 (time based uuid)
	u[8] &= 0x3F // clear variant
	u[8] |= 0x80 // set to IETF variant

	return u
}

func (u UUID) String() string {
	var offsets = [...]int{0, 2, 4, 6, 9, 11, 14, 16, 19, 21, 24, 26, 28, 30, 32, 34}
	const hexString = "0123456789abcdef"
	r := make([]byte, 36)
	for i, b := range u {
		r[offsets[i]] = hexString[b>>4]
		r[offsets[i]+1] = hexString[b&0xF]
	}
	r[8] = '-'
	r[13] = '-'
	r[18] = '-'
	r[23] = '-'
	return string(r)
}

func GetBigFile(filename string) {
	fi, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	for i := 0; i <= 21000; i++ {
		_, _ = fi.Write([]byte(GetRandomString(10000)))
	}
	fi.Write([]byte(";;;18899887788;;;330329199909092113"))
	fi.Close()
}

func SaveTemp(data map[string]interface{}) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&data)
	if err != nil {
		return err
	}
	file, err := os.OpenFile("./temp/"+strconv.FormatInt(time.Now().Unix(), 10)+".ini", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = file.Write(buffer.Bytes())
	if err != nil {
		return err
	}
	_ = file.Close()
	return nil
}

func ReadTemp(filename string) (map[string]interface{}, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 777)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	reader := bufio.NewReader(file)
	decoder := gob.NewDecoder(reader)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	_ = os.Remove(filename)
	_ = file.Close()
	return data, nil
}
