package voo

import "net/http"

func (v *Voo) SessionLoad(next http.Handler) http.Handler {
	return v.Session.LoadAndSave(next)
}
