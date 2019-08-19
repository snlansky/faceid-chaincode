package tpl

import (
	"reflect"
	"fmt"
	"strings"
	"unicode/utf8"
	"unicode"
)

type Parser struct {
	t       reflect.Type
	imports map[string]struct{}
	Fields  []string
	Methods []Method
}

func NewParser(v interface{}) *Parser {
	return &Parser{
		t:       reflect.TypeOf(v),
		imports: map[string]struct{}{},
		Fields:  []string{},
		Methods: []Method{},
	}
}

func (p *Parser) Parse() *Class {
	p.parseFields()
	p.parseMethod()
	return &Class{
		ClassName: p.splitTypeName(p.t.String()),
		Fields:    p.Fields,
		Imports:   p.imports,
		Methods:   p.Methods,
	}
}

func (p *Parser) splitTypeName(complete string) string {
	part := strings.Split(complete, ".")
	return part[len(part)-1]
}

func (p *Parser) parseFields() {
	var rawType reflect.Type
	if p.t.Kind() == reflect.Ptr {
		rawType = p.t.Elem()
	}
	for i := 0; i < rawType.NumField(); i++ {
		f := rawType.Field(i)
		char, _ := utf8.DecodeRuneInString(f.Name)
		if !unicode.IsUpper(char) {
			continue
		}
		var lastType string
		var lastName string
		t, find := javaType[f.Type.Kind()]
		if find {
			lastType = t
		} else {
			lastType = f.Type.Kind().String()
		}
		if v, find := iapk[f.Type.Kind()]; find {
			p.imports[v] = struct{}{}
		}

		if tranName := f.Tag.Get("json"); tranName == "" {
			lastName = f.Name
		} else {
			lastName = strings.Split(tranName, ",")[0]
		}
		p.Fields = append(p.Fields, fmt.Sprintf("private %s %s", lastType, lastName))

	}
}

func (p *Parser) parseMethod() {
	for i := 0; i < p.t.NumMethod(); i++ {
		m := p.t.Method(i)
		char, _ := utf8.DecodeRuneInString(m.Name)
		if !unicode.IsUpper(char) {
			continue
		}

		names, sig := p.parseParams(m)
		p.Methods = append(p.Methods, Method{FuncName: m.Name, ParamSig: sig, Params: names, RetType: p.parseRet(m)})
	}
}

func (p *Parser) parseParams(method reflect.Method) (varName []string, varSig []string) {

	l := method.Type.NumIn()
	if l < 2 {
		panic("method.Type.NumIn < 2")
	}

	for i := 2; i < l; i++ {
		t := method.Type.In(i)
		name := javaType[t.Kind()]
		if v, find := iapk[t.Kind()]; find {
			p.imports[v] = struct{}{}
		}
		v := fmt.Sprintf("%s%d", string(t.String()[0]), i-1)
		varSig = append(varSig, fmt.Sprintf("%s %s", name, v))
		varName = append(varName, v)
	}
	return
}

func (p *Parser) parseRet(method reflect.Method) string {

	l := method.Type.NumOut()
	if l == 0 {
		return "void"
	}
	if l > 1 {
		panic("method.Type.NumOut != 1")
	}
	t := method.Type.Out(0)
	s := strings.Split(t.String(), ".")
	name := s[len(s)-1]
	javaName, find := javaType[t.Kind()]
	if find {
		return javaName
	}
	if v, find := iapk[t.Kind()]; find {
		p.imports[v] = struct{}{}
	}
	return name
}
