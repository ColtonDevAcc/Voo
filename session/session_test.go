package session

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/alexedwards/scs/v2"
)

func TestSession_InitSession(t *testing.T) {

	v := &Session{
		CookieLifetime: "100",
		CookiePersist:  "true",
		CookieName:     "voo",
		CookieDomain:   "localhost",
		SessionType:    "cookie",
	}

	var sm *scs.SessionManager

	ses := v.InitSession()

	var sessKind reflect.Kind
	var sessType reflect.Type

	rv := reflect.ValueOf(ses)

	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		fmt.Println("for loop:", rv.Kind(), rv.Kind(), rv)

		sessKind = rv.Kind()
		sessType = rv.Type()

		rv = rv.Elem()
	}

	if !rv.IsValid() {
		t.Error("invalid type or king; kind: ", rv.Kind(), " type: ", rv.Type())
	}

	if sessKind != reflect.ValueOf(sm).Kind() {
		t.Error("wrong kind return testing cookie session. Expected ", reflect.ValueOf(sm).Kind(), " got ", sessKind)
	}

	if sessType != reflect.ValueOf(sm).Type() {
		t.Error("wrong type return testing cookie session. Expected ", reflect.ValueOf(sm).Type(), " got ", sessType)
	}
}
