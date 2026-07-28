package controllers

import "net/http"

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("login successful"))
}
