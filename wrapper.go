package wrapper

import (
	"xorm.io/xorm"
)

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
