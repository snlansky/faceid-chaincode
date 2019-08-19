package rpc

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"reflect"
	"fmt"
	"testing"
)

func Test_convertStringToParam(t *testing.T) {

	s := "lilith"
	v, ok := convertString2Value(s, reflect.TypeOf("name"))
	assert.True(t, ok)
	assert.Equal(t, v.Interface(), s)

	s = "12"
	v, ok = convertString2Value(s, reflect.TypeOf(100))
	assert.True(t, ok)
	assert.Equal(t, v.Interface(), 12)

	s = "12"
	v, ok = convertString2Value(s, reflect.TypeOf(uint8(0)))
	assert.True(t, ok)
	assert.Equal(t, v.Interface(), uint8(12))

	s = "12"
	v, ok = convertString2Value(s, reflect.TypeOf(uint16(0)))
	assert.True(t, ok)
	assert.Equal(t, v.Interface(), uint16(12))

	s = "12"
	v, ok = convertString2Value(s, reflect.TypeOf(uint32(0)))
	assert.True(t, ok)
	assert.Equal(t, v.Interface(), uint32(12))

	s = "12"
	v, ok = convertString2Value(s, reflect.TypeOf(uint64(0)))
	assert.True(t, ok)
	assert.Equal(t, v.Interface(), uint64(12))

	s = "12"
	v, ok = convertString2Value(s, reflect.TypeOf(int8(0)))
	assert.True(t, ok)
	assert.Equal(t, v.Interface(), int8(12))

	s = "12"
	v, ok = convertString2Value(s, reflect.TypeOf(int16(0)))
	assert.True(t, ok)
	assert.Equal(t, v.Interface(), int16(12))

	s = "12"
	v, ok = convertString2Value(s, reflect.TypeOf(int32(0)))
	assert.True(t, ok)
	assert.Equal(t, v.Interface(), int32(12))

	s = "12"
	v, ok = convertString2Value(s, reflect.TypeOf(int64(0)))
	assert.True(t, ok)
	assert.Equal(t, v.Interface(), int64(12))

	s = "3.1415"
	v, ok = convertString2Value(s, reflect.TypeOf(float32(0.0)))
	assert.True(t, ok)
	assert.Equal(t, v.Interface(), float32(3.1415))

	s = "12.12345"
	v, ok = convertString2Value(s, reflect.TypeOf(float64(0.0)))
	assert.True(t, ok)
	assert.Equal(t, v.Interface(), float64(12.12345))

	s = "12.123"
	_, ok = convertString2Value(s, reflect.TypeOf(int32(0)))
	assert.False(t, ok)

	s = "12"
	v, ok = convertString2Value(s, reflect.TypeOf(float64(0.0)))
	assert.True(t, ok)
	assert.Equal(t, v.Interface(), float64(12))

}

func Test_convertParam1(t *testing.T) {
	var p interface{}
	p = map[string]interface{}{"s1": 56, "s2": 567.89, "s3": true, "s4": "v1"}
	v, ok := convertParam(p, reflect.TypeOf(map[string]interface{}{}))
	assert.True(t, ok)
	fmt.Println("v:", v.Interface())
	assert.True(t, reflect.DeepEqual(v.Interface(), p))

	p = map[string]interface{}{"k1": "s3", "k2": "s5"}
	v, ok = convertParam(p, reflect.TypeOf(map[string]string{}))
	assert.True(t, ok)
	fmt.Println(v.Interface())
	p1 := p.(map[string]interface{})
	v1 := v.Interface().(map[string]string)
	assert.Equal(t, len(p1), len(v1))
	for k, v := range p1 {
		assert.Equal(t, v.(string), v1[k])
	}

	p = map[string]interface{}{"k1": 12, "k2": 5}
	v, ok = convertParam(p, reflect.TypeOf(map[string]int{}))
	assert.True(t, ok)
	fmt.Println(v.Interface())
	p2 := p.(map[string]interface{})
	v2 := v.Interface().(map[string]int)
	assert.Equal(t, len(p2), len(v2))
	for k, v := range p2 {
		assert.Equal(t, v.(int), v2[k])
	}

	p = map[string]interface{}{"k1": false, "k2": true}
	v, ok = convertParam(p, reflect.TypeOf(map[string]bool{}))
	assert.True(t, ok)
	fmt.Println(v.Interface())
	p3 := p.(map[string]interface{})
	v3 := v.Interface().(map[string]bool)
	assert.Equal(t, len(p3), len(v3))
	for k, v := range p3 {
		assert.Equal(t, v.(bool), v3[k])
	}
	//
	p = map[string]interface{}{"k1": 344.78, "k2": 544.89}
	v, ok = convertParam(p, reflect.TypeOf(map[string]float64{}))
	assert.True(t, ok)
	fmt.Println(v.Interface())
	p4 := p.(map[string]interface{})
	v4 := v.Interface().(map[string]float64)
	assert.Equal(t, len(p4), len(v4))
	for k, v := range p4 {
		assert.Equal(t, v.(float64), v4[k])
	}

	p = []interface{}{"s1", 34, 554.99, true}
	v, ok = convertParam(p, reflect.TypeOf([]interface{}{}))
	assert.True(t, ok)
	fmt.Println(v.Interface())
	assert.True(t, reflect.DeepEqual(v.Interface(), p))

	p = []interface{}{"s1", "s2"}
	v, ok = convertParam(p, reflect.TypeOf([]string{}))
	assert.True(t, ok)
	fmt.Println(v)
	p5 := p.([]interface{})
	v5 := v.Interface().([]string)
	assert.Equal(t, len(p5), len(v5))
	for k, v := range p5 {
		assert.Equal(t, v.(string), v5[k])
	}

	p = []interface{}{34.67, 345.89}
	v, ok = convertParam(p, reflect.TypeOf([]float64{}))
	assert.True(t, ok)
	fmt.Println(v)
	p6 := p.([]interface{})
	v6 := v.Interface().([]float64)
	assert.Equal(t, len(p6), len(v6))
	for k, v := range p6 {
		assert.Equal(t, v.(float64), v6[k])
	}

	p = []interface{}{34.67, 345.89}
	v, ok = convertParam(p, reflect.TypeOf([]float32{}))
	assert.True(t, ok)
	fmt.Println(v)
	p7 := p.([]interface{})
	v7 := v.Interface().([]float32)
	assert.Equal(t, len(p7), len(v7))
	for k, v := range p7 {
		assert.Equal(t, float32(v.(float64)), v7[k])
	}

	p = []interface{}{true, false}
	v, ok = convertParam(p, reflect.TypeOf([]bool{}))
	assert.True(t, ok)
	fmt.Println(v)
	p8 := p.([]interface{})
	v8 := v.Interface().([]bool)
	assert.Equal(t, len(p8), len(v8))
	for k, v := range p8 {
		assert.Equal(t, v.(bool), v8[k])
	}
}

func Test_convertParam(t *testing.T) {
	type Person struct {
		Name      string  `json:"name"`
		Age       int     `json:"age,omitempty"`
		Money     float64 `json:"money,omitempty"`
		IsStudent bool    `json:"is_student,omitempty"`
	}
	type Work struct {
		Time   string
		Num    int
		Safe   bool
		Salary float64
	}
	type Complex struct {
		Name   string  `json:"name"`
		Member float64 `json:"member"`
		P      *Person
		W      Work
		Mi     map[string]interface{}
		Ms     map[string]string
		Seti   []interface{}
		SetN   []float64
	}
	p := &Person{
		Name:      "lucy",
		Age:       21,
		Money:     123.46,
		IsStudent: false,
	}

	var params []interface{}
	s := `[12,true,123.567,{"name":"lucy","age":21,"money":123.46,"is_student":false}]`
	err := json.Unmarshal([]byte(s), &params)
	assert.NoError(t, err)
	v, ok := convertParam(params[3], reflect.TypeOf(p))
	assert.True(t, ok)
	fmt.Println(v.Interface())
	assert.True(t, reflect.DeepEqual(v.Interface(), p))

	c := &Complex{
		Name:   "1class",
		Member: 45,
		P:      p,
		W: Work{
			Time:   "12月15日",
			Num:    12,
			Safe:   true,
			Salary: 1234.89,
		},
		Mi:   map[string]interface{}{"s1": 56, "s2": 567.89, "s3": true, "s4": "v1"},
		Ms:   map[string]string{"k1": "s3", "k2": "s5"},
		Seti: []interface{}{89, true, 45.97, "v2"},
		SetN: []float64{1234.66, 768.89},
	}
	classStr, _ := json.Marshal(c)
	var arg interface{}
	err = json.Unmarshal([]byte(classStr), &arg)
	assert.NoError(t, err)
	v, ok = convertParam(arg, reflect.TypeOf(c))
	assert.True(t, ok)
	fmt.Println(v.Interface())
	fmt.Println(c)
}
