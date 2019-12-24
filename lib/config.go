package lib

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var Config = make(map[string]interface{}, 1024)

func init() {
	Config["Temp"] = ""
	file, err := os.OpenFile("./config.ini", os.O_RDONLY, 777)

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	var lineNo int
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
		lineNo++
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '\n' || line[0] == '#' {
			continue
		}

		arr := strings.Split(line, "=")
		if len(arr) == 0 {
			fmt.Println("Invalid Config,Line:%d\n", lineNo)
			continue
		}
		key := strings.TrimSpace(arr[0])
		if len(key) == 0 {
			fmt.Println("Invalid Config,Line:%d\n", lineNo)
			continue
		}
		if len(arr) == 1 {
			Config[key] = ""
			continue
		}
		value := strings.TrimSpace(strings.Trim(arr[1], "\n"))
		Config[key] = value
	}
	file.Close()
}
