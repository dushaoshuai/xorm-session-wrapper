package main

import (
	"os"
	"reflect"
	"text/template"

	"xorm.io/xorm"
)

//go:generate bash ./gen.sh

func main() {
	gen(methods())
}

type in struct {
	Type       string
	IsVariadic bool
}

type method struct {
	Name string
	Ins  []in
}

func methods() (methods []method) {
	sessType := reflect.TypeOf((*xorm.Session)(nil))

	for i := 0; i < sessType.NumMethod(); i++ {
		methodType := sessType.Method(i).Type

		if methodType.NumOut() != 1 {
			continue
		}
		if methodType.Out(0) != sessType {
			continue
		}

		m := method{
			Name: sessType.Method(i).Name,
			Ins:  nil,
		}

		numIn := methodType.NumIn()
		if numIn == 0 { // impossible
			continue
		}

		for j := 1; j < numIn; j++ {
			inType := methodType.In(j)

			if j == numIn-1 && methodType.IsVariadic() {
				m.Ins = append(m.Ins, in{
					Type:       "..." + inType.Elem().String(),
					IsVariadic: true,
				})
			} else {
				m.Ins = append(m.Ins, in{
					Type:       inType.String(),
					IsVariadic: false,
				})
			}
		}

		methods = append(methods, m)
	}

	return
}

func subtract(a, b int) int {
	return a - b
}

func gen(data any) {
	fnMap := template.FuncMap{
		"subtract": subtract,
	}

	tmplvar := template.Must(
		template.New("tmpl").Funcs(fnMap).Parse(tmpl),
	)

	f, err := os.Create("../session.go")
	if err != nil {
		panic(err)
	}

	err = tmplvar.Execute(f, data)
	if err != nil {
		panic(err)
	}
}
