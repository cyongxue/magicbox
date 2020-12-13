package xhvalidate

import (
	myValidator "github.com/golang-collections/lib.go/validation/validator"
	"github.com/kataras/iris/v12"
	"gopkg.in/go-playground/validator.v9"
)

var Validate *validator.Validate

// ValidateMobile 自定义验证电话号码
func ValidateCNMobile(fl validator.FieldLevel) bool {
	return myValidator.IsCnMobile(fl.Field().String())
}

func init() {
	Validate = validator.New()
	// 自定义增加tag  is_CNMobile
	_ = Validate.RegisterValidation("is_CNMobile", ValidateCNMobile)
}

func CheckCNMobile(mobile string) bool {
	return myValidator.IsCnMobile(mobile)
}

func CheckCNPhone(phone string) bool {
	return myValidator.IsCnPhone(phone)
}

// JSONInputValidate json post提交的数据进行解析处理
func JSONInputValidate(ctx iris.Context, jsonObjPtr interface{}) error {
	if err := ctx.ReadJSON(jsonObjPtr); err != nil {
		return err
	}
	if err := Validate.Struct(jsonObjPtr); err != nil {
		return err
	}
	return nil
}
