package main

var tmpl = `// Code generated by ./gen. DO NOT EDIT.

package wrapper

import (
	"context"

	"xorm.io/xorm"
	"xorm.io/xorm/contexts"
)

type session struct {
	*xorm.Session
}

{{range $ -}}

func (s *session) {{.Name}}(

{{- $l := len .Ins -}}

{{- range $i, $e := .Ins -}}

x{{$i}} {{$e.Type}}{{if lt $i (subtract $l 1)}}, {{end}}

{{- end -}}

) *Session {
	s.Session.{{.Name}}(

{{- range $i, $e := .Ins -}}

x{{$i}}{{if lt $i (subtract $l 1)}}, {{end}}

{{- end -}}

)
	return &Session{s}
}

{{end}}
`
