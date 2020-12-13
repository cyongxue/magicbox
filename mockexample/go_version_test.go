package mockexample

import (
	"github.com/cyongxue/magicbox/mockexample/spider"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestGetGoVersion(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockSpider := spider.NewMockSpider(mockCtl)
	mockSpider.EXPECT().GetBody().Return("go1.8.3")

	goVer := GetGoVersion(mockSpider)
	if goVer != "go1.8.3" {
		t.Errorf("Get wrong version %s", goVer)
	}

	mockSpider.EXPECT().GetBody().Return("1.9.1")
	goVer = GetGoVersion(mockSpider)
	if goVer != "go1.8.3" {
		t.Errorf("Get wrong version %s", goVer)
	}

}
