package utils

type ErrorJsonStruct struct {
	oper_code int
}

var JsonError ErrorJsonStruct

func init() {
	JsonError := &ErrorJsonStruct{}
	JsonError.oper_code = 0
}
