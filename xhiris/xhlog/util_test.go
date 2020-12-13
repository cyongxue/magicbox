package xhlog

import (
	"fmt"
	"testing"
	"time"
)

func TestNowMsec(t *testing.T) {
	fmt.Println(time.Now().Unix())
	fmt.Println(NowMsec())
	fmt.Println(time.Now().UnixNano())
}
