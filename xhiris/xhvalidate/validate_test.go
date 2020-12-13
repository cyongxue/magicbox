package xhvalidate

import (
	"fmt"
	"testing"
)

func TestValidateCNMobile(t *testing.T) {
	type Input struct {
		Mobile   string `json:"mobile" validate:"required,is_CNMobile"`
		Password string `json:"password" validate:"required"`
		Nick     string `json:"nick" validate:"required"`
	}

	fmt.Println(CheckCNMobile("132456786541"))
}
