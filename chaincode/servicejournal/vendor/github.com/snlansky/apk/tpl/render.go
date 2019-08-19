package tpl

import (
	"reflect"
	"text/template"
	"os"
)

var iapk = map[reflect.Kind]string{
	reflect.Map: "import java.util.HashMap",
}
var javaType = map[reflect.Kind]string{
	reflect.Int:     "int",
	reflect.Map:     "HashMap",
	reflect.String:  "String",
	reflect.Int64:   "Long",
	reflect.Bool:    "Boolean",
	reflect.Float64: "float",
}

func Render(fname func(name string)string, v interface{}, tplFile string) error {

	p := NewParser(v)
	c := p.Parse()

	funcMap := template.FuncMap{}
	tpl, err := template.ParseFiles(tplFile)
	if err != nil {
		return err
	}
	tpl.Funcs(funcMap)

	f, err := os.OpenFile(fname(c.ClassName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
	if err != nil {
		return err
	}

	return tpl.Execute(f, c)

}
