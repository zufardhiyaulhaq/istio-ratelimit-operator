package utils

import (
	"github.com/champly/lib4go/encoding"
	proto_types "github.com/gogo/protobuf/types"
)

func ConvertYaml2Struct(str string) *proto_types.Struct {
	res, _ := encoding.YAML2Struct(str)
	return res
}
