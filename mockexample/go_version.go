package mockexample

import "github.com/cyongxue/magicbox/mockexample/spider"

func GetGoVersion(s spider.Spider) string {
	body := s.GetBody()
	return body
}
