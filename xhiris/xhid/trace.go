package xhid

import (
	"crypto/md5"
	"fmt"
	"io"
	"sync"
	"time"
)

const (
	TraceId  = "x-traceid"
	SpanId   = "x-spanid"
	ParentId = "x-parentid"
)

/**
uuid的设计原则：
   1. 唯一性；
   2. 时间相关；
   3. 粗略有序；
   4. 可反解；
   5. 可制造；

用于int64类型的uuid生成。
目前的用处：
   1. weapon id的生成；

采用秒(s)级weibo的实现方式实现
设计细节：
   预留(64 - 50)位用于后期的处理
   30bit       秒级时间，timestamp
   16bit       序列号，sequence
   4bit        区分idc，目前
*/
type SnowFlakeId struct {
	Idc          int64 // idc区分，wolverine中可以用于区分域
	Sequence     int64
	SecTimeStamp int64 // s秒单位时间戳
}

var once sync.Once
var idDriver *SnowFlakeId

func IdDriver(idc int64) *SnowFlakeId {
	once.Do(func() {
		idDriver = &SnowFlakeId{Idc: int64(idc), Sequence: 0, SecTimeStamp: 0}
	})
	return idDriver
}

func (s *SnowFlakeId) GetNextId() int64 {
	timeStamp := time.Now().Unix()

	if s.SecTimeStamp != timeStamp { // 时间戳不相等
		s.Sequence = 0
	} else { // 时间戳相等，递增sequence
		s.Sequence = (s.Sequence + 1) & 0xFFFF
	}

	s.SecTimeStamp = timeStamp

	var id int64
	id = (s.Sequence & 65535) | ((s.Idc & 15) << 16) | ((s.SecTimeStamp & 1073741823) << 20)
	return id
}

func MakeSpanId(origin string) string {
	w := md5.New()
	io.WriteString(w, fmt.Sprintf("%s%d", origin, time.Now().UnixNano()))
	md5str := fmt.Sprintf("%x", w.Sum(nil))
	return md5str[0:15]
}
