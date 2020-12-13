package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	TimeFormat = "20060102_1504"
)

type LogRotateConfig struct {
	LogPrefix  string `json:"log_prefix"`
	LogPath    string `json:"path"`
	RotatePath string `json:"rotate_path"`
	BeforeHour int    `json:"before_hour"`
	DeleteHour int    `json:"delete_hour"`
}

func main() {
	runtime.GOMAXPROCS(1)

	cfg := flag.String("json", "{}", "配置文件: {\"log_prefix\": \"server_\", "+
		"\"path\": \"/Users/hehui/GolandProjects/src/github.com/cyongxue/magicbox/xhiris/xhlogrotate\","+
		"\"rotate_path\": \"/Users/hehui/GolandProjects/src/github.com/cyongxue/magicbox/xhiris/xhlogrotate/rotate\",\"before_hour\": 8}")
	flag.Parse()

	config := LogRotateConfig{}
	if err := json.Unmarshal([]byte(*cfg), &config); err != nil {
		fmt.Println(err)
		return
	}
	//config = LogRotateConfig{
	//	LogPrefix:  "server_",
	//	LogPath:    "/Users/hehui/GolandProjects/src/github.com/cyongxue/magicbox/xhiris/xhlogrotate",
	//	RotatePath: "/Users/hehui/GolandProjects/src/github.com/cyongxue/magicbox/xhiris/xhlogrotate/logs",
	//	BeforeHour: 8,
	//	DeleteHour: 72,
	//}
	fmt.Println(config)

	// 压缩文件
	if err := walFile(&config); err != nil {
		fmt.Println(err)
	}

	// 删除存档文件
	tryDeleteRotateLog(&config)
	return
}

func tryDeleteRotateLog(config *LogRotateConfig) {
	_ = filepath.Walk(config.RotatePath, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasPrefix(info.Name(), config.LogPrefix) {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".gz") {
			return nil
		}

		name := info.Name()
		timeFormat := name[strings.Index(name, "_")+1 : strings.Index(name, ".")]
		if len(timeFormat) > len(TimeFormat) {
			timeFormat = timeFormat[:len(TimeFormat)]
		}

		t, err := time.Parse(TimeFormat, timeFormat)
		if err != nil {
			return err
		}
		now := time.Now().UTC()
		if (now.Unix() - t.Unix()) > int64(config.DeleteHour*3600) {
			os.Remove(path)
		}

		return nil
	})

	return
}

func walFile(config *LogRotateConfig) error {
	err := filepath.Walk(config.LogPath, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasPrefix(info.Name(), config.LogPrefix) {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".log") {
			return nil
		}

		name := info.Name()
		timeFormat := name[strings.Index(name, "_")+1 : strings.Index(name, ".")]
		if len(timeFormat) > len(TimeFormat) {
			timeFormat = timeFormat[:len(TimeFormat)]
		}

		t, err := time.Parse(TimeFormat, timeFormat)
		if err != nil {
			return err
		}
		now := time.Now().UTC()
		if (now.Unix() - t.Unix()) > int64(config.BeforeHour*3600) {
			// 文件压缩归档
			destFileName := strings.Split(info.Name(), ".")[0] + ".gz"
			if err := compress(config, destFileName, path); err != nil {
				fmt.Println(err)
			} else {
				os.Remove(path)
				os.Rename(destFileName, filepath.Join(config.RotatePath, destFileName))
			}
		}

		return nil
	})
	return err
}

func compress(config *LogRotateConfig, dest string, origin string) error {
	outputFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	gzipWriter := gzip.NewWriter(outputFile)
	defer gzipWriter.Close()

	inputFile, err := os.Open(origin)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	reader := bufio.NewReader(inputFile)
	for {
		s, e := reader.ReadString('\n')
		if e == io.EOF {
			break
		}
		_, err := gzipWriter.Write([]byte(s))
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
