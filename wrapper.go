package wrapper

import (
	"fmt"
	"reflect"
	"strings"

	"xorm.io/xorm"
)

// Session is a thin wrapper of *xorm.Session.
// It overrides some methods of *xorm.Session and
// delegates others to the embedding *xorm.Session.
// It aims at eliminating tedious if statements.
type Session struct {
	*session
}

// NewSession returns a new Session that wraps the embedded *xorm.Session.
func NewSession(embedded *xorm.Session) *Session {
	return &Session{
		session: &session{
			Session: embedded,
		},
	}
}

// In overrides (*xorm.Session).In method. If no values are given,
// or if the first value is nil, or if the first value is an empty slice,
// it does nothing. Otherwise, it delegates to (*xorm.Session).In.
func (s *Session) In(column string, values ...any) *Session {
	if len(values) <= 0 {
		return s
	}

	firstVal := values[0]
	if firstVal == nil {
		return s
	}

	rv := reflect.ValueOf(firstVal)
	if !rv.IsValid() {
		return s
	}

	if rv.Kind() != reflect.Slice {
		s.Session.In(column, values...)
		return s
	}

	if rv.Len() > 0 {
		s.Session.In(column, values...)
		return s
	}

	return s
}

// Equal builds a `column = val` condition if val is not nil and is not the zero value.
func (s *Session) Equal(column string, val any) *Session {
	if val == nil {
		return s
	}

	rt := reflect.TypeOf(val)
	if val == reflect.Zero(rt).Interface() {
		return s
	}

	s.Session.Where(
		fmt.Sprintf("%s = ?", column), val,
	)
	return s
}

// Ranger defines a range with a Start and an End.
type Ranger struct {
	Start any
	End   any
}

// Between builds a `column BETWEEN ranger.Start AND ranger.End` condition if ranger is not nil.
func (s *Session) Between(column string, ranger *Ranger) *Session {
	if ranger == nil {
		return s
	}

	s.Session.Where(
		fmt.Sprintf("%s BETWEEN ? AND ?", column),
		ranger.Start, ranger.End,
	)
	return s
}

// Like builds a `column LIKE %val%` condition with val (strings.TrimSpace)ed.
// The condition is built only if the trimmed val is not an empty string.
func (s *Session) Like(column string, val string) *Session {
	val = strings.TrimSpace(val)
	if val == "" {
		return s
	}

	// escape wildcard characters
	val = strings.ReplaceAll(val, `_`, `\_`)
	val = strings.ReplaceAll(val, `%`, `\%`)

	s.Session.Where(
		fmt.Sprintf("%s LIKE ?", column),
		"%"+val+"%",
	)
	return s
}
