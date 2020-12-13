package xhlog

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	// 日志文件后缀时间格式
	suffixFormat      = "20060102_1504" // 切分时，记录到分钟级别
	DefaultRotateSize = 100 * 1024 * 1024
)

// 日志文件方式输出
// 1. 全部日志输出，包括各个级别的日志
// 2. 全量日志是默认，prefix.log的方式
type fileLogger struct {
	cancel     context.CancelFunc
	fileDir    string // 日志文件目录
	prefix     string // 日志文件的前缀，写日志时的文件名
	rotateSize int64  // M为单位，缺省采用100M

	logMux      sync.RWMutex
	logFile     *os.File // 写日志的文件
	logChan     chan string
	currentSize int64 // 当前文件的大小
}

// newFileLogger 创建一个新的 fileLogger
func newFileLogger(conf *LoggerConf) (*fileLogger, error) {
	f := &fileLogger{
		fileDir:    conf.Dir,
		prefix:     conf.Prefix,
		rotateSize: conf.RotateSize,
		logChan:    make(chan string, 5000),
	}
	var ctx context.Context
	ctx, f.cancel = context.WithCancel(context.Background())

	f.isExistOrCreate()

	var err error
	logFilePrefix := filepath.Join(f.fileDir, f.prefix)
	logFileName := fmt.Sprintf("%s.log", logFilePrefix)
	// 尝试将现有的log补偿切分
	if fileExists(logFileName) {
		logFileRotate := fmt.Sprintf("%s_%s.log", logFilePrefix, time.Now().Format(suffixFormat))
		_ = os.Rename(logFileName, logFileRotate)
	}

	f.logFile, err = os.OpenFile(logFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("open logFile error: %s", err.Error())
	}

	go f.monitor(ctx)
	go f.flush(ctx)

	return f, nil
}

// fileExists 判断文件是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// fileLogger.close 关闭日志打印，进行资源的回收
func (f *fileLogger) close() {
	f.cancel()
	time.Sleep(1 * time.Second)

	f.logFile.Close()
	close(f.logChan)
	return
}

// fileLogger.isExistOrCreate 尝试创建日志文件目录
func (f *fileLogger) isExistOrCreate() {
	_, err := os.Stat(f.fileDir)
	if err != nil && os.IsNotExist(err) {
		_ = os.Mkdir(f.fileDir, 0755)
	}
}

// fileLogger.split 执行日志切割
func (f *fileLogger) split() error {
	f.logMux.Lock()
	defer f.logMux.Unlock()

	logFilePrefix := filepath.Join(f.fileDir, f.prefix)
	logFileName := fmt.Sprintf("%s.log", logFilePrefix)
	logFileRotate := fmt.Sprintf("%s_%s.log", logFilePrefix, time.Now().Format(suffixFormat))

	if f.logFile != nil {
		f.logFile.Close()
	}
	if err := os.Rename(logFileName, logFileRotate); err != nil {
		f.logFile, _ = os.OpenFile(logFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
		return fmt.Errorf("rename logfile to logfileRotate error: %s", err.Error())
	}

	var err error
	f.logFile, err = os.OpenFile(logFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("open logFile error: %s", err.Error())
	}
	f.currentSize = 0

	return nil
}

// fileLogger.isMustSplit 判断是否到了切割点
func (f *fileLogger) isMustSplit() bool {
	if f.currentSize >= f.rotateSize {
		return true
	}
	return false
}

// fileLogger.monitor 异步监听，进行日志切割
func (f *fileLogger) monitor(ctx context.Context) {
	defer func() {
		recover()
	}()

	t := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("fileLogger monitor return")
			return
		case <-t.C:
			if f.isMustSplit() {
				if err := f.split(); err != nil {
					// todo: 错误日志打印
				}
			}
		}
	}
}

// fileLogger.flush 日志下刷
func (f *fileLogger) flush(ctx context.Context) {
	defer func() {
		recover()
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("fileLogger flush return\n")
			return
		case str := <-f.logChan:
			if f.logFile != nil {
				f.logMux.RLock()
				f.logFile.Write([]byte(str))
				f.currentSize = f.currentSize + int64(len(str))
				f.logMux.RUnlock()
			}
		}
	}
}
