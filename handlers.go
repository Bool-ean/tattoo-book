package main

import (
	"fmt"
	"log"
	"net/http"
)

// function for handling http request at index "/"
func serveIndex(w http.ResponseWriter, r *http.Request) {
	//honestly not sure why this error check would be here
	//how would this function get called we we weren't at "/"
	if r.URL.Path != "/" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	//maing sure request is GET
	if r.Method != "GET" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, "templates/index.html")
}


func serveLogin(hub *Hub, w http.ResponseWriter, r *http.Request) {
	//honestly not sure why this error check would be here
	//how would this function get called we we weren't at "/"
	if r.URL.Path != "/login" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	switch r.Method{

	case http.MethodGet:
		http.ServeFile(w, r, "templates/login.html")

	case http.MethodPost:
		log.Println("inside post")
		un := r.FormValue("un")
		//pw := r.FormValue("pw")

		//check if username exists
		if hub.UsernameTaken(un){
			log.Println("inside username take true")
			//if it doesn't, return error
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `<div id="login-error" class="mt-4 text-red-600 font-medium">Username already taken</div>`)
			return
		}
		//if it does, return succes, redirect to /
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

	}
}

