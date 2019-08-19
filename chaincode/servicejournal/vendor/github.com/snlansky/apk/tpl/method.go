package tpl

import (
	"fmt"
	"strings"
)

type Class struct {
	ClassName string
	Fields    []string
	Imports   map[string]struct{}
	Methods   []Method
}
type Method struct {
	FuncName string
	ParamSig []string
	Params   []string
	RetType  string
}

func (m Method) GetSig() string {
	return fmt.Sprintf("public %s %s(%s)", m.RetType, m.FuncName, strings.Join(m.ParamSig, ", "))
}

func (m Method) GenInterface(s string) string {
	return fmt.Sprintf("%s.%s", s, m.FuncName)
}

func (m Method) GetParams() string {
	return fmt.Sprintf("{%s}", strings.Join(m.Params, ","))
}
