package voo

import "net/http"

func (v *Voo) SessionLoad(next http.Handler) http.Handler {
	v.InfoLog.Println("SessionLoad called")

	return v.Session.LoadAndSave(next)
}
