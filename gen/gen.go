package main

import (
	_ "embed"
	"os"
	"reflect"
	"text/template"

	"xorm.io/xorm"
)

//go:generate bash ./gen.sh

//go:embed template.tmpl
var tmpl string

func main() {
	gen(getMethods())
}

type in struct {
	Type string
}

type method struct {
	Name       string
	Ins        []in
	IsVariadic bool
}

func getMethods() (methods []method) {
	sessType := reflect.TypeOf((*xorm.Session)(nil))

	for i := 0; i < sessType.NumMethod(); i++ {
		methodType := sessType.Method(i).Type

		// overrides only methods that return exclusively a *xorm.Session
		if methodType.NumOut() != 1 {
			continue
		}
		if methodType.Out(0) != sessType {
			continue
		}

		m := method{
			Name:       sessType.Method(i).Name,
			Ins:        nil,
			IsVariadic: methodType.IsVariadic(),
		}

		numIn := methodType.NumIn()
		if numIn == 0 { // impossible
			continue
		}

		for j := 1; j < numIn; j++ {
			inType := methodType.In(j)

			if j == numIn-1 && methodType.IsVariadic() {
				m.Ins = append(m.Ins, in{
					Type: "..." + inType.Elem().String(),
				})
			} else {
				m.Ins = append(m.Ins, in{
					Type: inType.String(),
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
