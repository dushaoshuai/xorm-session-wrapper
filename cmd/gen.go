package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"xorm.io/xorm"
)

func main() {
	var s *xorm.Session
	st := reflect.TypeOf(s)

	for i := 0; i < st.NumMethod(); i++ {
		method := st.Method(i)
		name := method.Name
		typ := method.Type

		if typ.NumOut() != 1 {
			continue
		}
		if typ.Out(0) != st {
			continue
		}

		numIn := typ.NumIn()
		if numIn == 0 {
			continue
		}

		var builder strings.Builder
		builder.WriteString(name)
		builder.WriteByte('(')

		for j := 1; j < numIn; j++ {
			builder.WriteByte('x')
			builder.WriteString(strconv.Itoa(j))
			builder.WriteByte(' ')

			inType := typ.In(j)

			if j == numIn-1 && typ.IsVariadic() {
				builder.WriteString("...")
				builder.WriteString(inType.Elem().String())
			} else {
				builder.WriteString(inType.String())
			}

			if j < numIn-1 {
				builder.WriteString(", ")
			}
		}

		builder.WriteString(") *session")

		fmt.Println(builder.String())
	}
}
