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
type Session struct {
	*session
}

func NewSession(embedded *xorm.Session) *Session {
	return &Session{
		session: &session{
			Session: embedded,
		},
	}
}

// In overrides (*xorm.Session).In method.
func (s *Session) In(column string, values ...any) *Session {
	if len(values) <= 0 {
		return s
	}

	rv := reflect.ValueOf(values[0])
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

type Ranger struct {
	Start any
	End   any
}

func (s *Session) Between(column string, ranger *Ranger) *Session {
	if ranger == nil {
		return s
	}

	s.Session.Where(
		fmt.Sprintf("%s between ? and ?", column),
		ranger.Start, ranger.End,
	)
	return s
}

func (s *Session) Like(column string, val string) *Session {
	val = strings.TrimSpace(val)
	if val == "" {
		return s
	}
	s.Session.Where(
		fmt.Sprintf("%s LIKE ?", column),
		"%"+val+"%",
	)
	return s
}
