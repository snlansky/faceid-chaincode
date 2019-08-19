package tpl

import (
	"testing"
	"fmt"
)

func TestRender(t *testing.T) {
	f := func(name string) string {
		return fmt.Sprintf("%s.java", name)
	}
	err := Render(f, nil, "class.tpl")
	if err != nil {
		panic(err)
	}
}
