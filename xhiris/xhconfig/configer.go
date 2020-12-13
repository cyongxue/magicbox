package xhconfig

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// !!!!!!!! 实现该接口的对象采用指针方式 !!!!!!!!!!!!!!!!!
type Parser interface {
	Parse(container ConfContainer, runMode string) error
}

// ConfigParse 配置文件解析的总体入口
func ConfigParse(parser interface{}, fileName string) error {
	conf := IniConfig{}
	container, err := conf.Parse(fileName)
	if err != nil {
		return err
	}
	runMode := container.String("runmode", nil)

	object := reflect.ValueOf(parser)
	refElem := object.Elem()
	typeOfType := object.Type()
	for i := 0; i < refElem.NumField(); i++ {
		field := refElem.Field(i)
		if field.IsNil() {
			return fmt.Errorf("field is nil, element index is %d", i)
		}
		configParser, ok := field.Interface().(Parser)
		if !ok {
			return fmt.Errorf("field no ConfigParser, name is %s", typeOfType.Field(i).Name)
		}
		if err := configParser.Parse(container, runMode); err != nil {
			return err
		}
	}

	data, err := json.Marshal(parser)
	if err != nil {
		return err
	}
	fmt.Printf("ServerConfig Parse: %s\n", string(data))
	return nil
}
