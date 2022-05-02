package session

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
)

type Session struct {
	CookieLifetime string
	CookiePersist  string
	CookieName     string
	CookieDomain   string
	SessionType    string
	CookieSecure   string
}

func (v *Session) InitSession() *scs.SessionManager {
	var persist, secure bool

	//? how long should sessions last?
	minutes, err := strconv.Atoi(v.CookieLifetime)
	if err != nil {
		minutes = 60
	}

	//! should cookies persist
	if strings.ToLower(v.CookiePersist) == "true" {
		persist = true
	}

	//! should cookies be secure
	if strings.ToLower(v.CookieSecure) == "true" {
		secure = true
	} else {
		secure = false
	}

	//create the session
	session := scs.New()
	session.Lifetime = time.Duration(minutes) * time.Minute
	session.Cookie.Persist = persist
	session.Cookie.Name = v.CookieName
	session.Cookie.Secure = secure
	session.Cookie.Domain = v.CookieDomain
	session.Cookie.SameSite = http.SameSiteLaxMode

	switch v.SessionType {
	case "redis":
	case "mysql", "mariadb":
	case "postgres", "postgresql":
	default:
		// cookie
	}

	return session
}
