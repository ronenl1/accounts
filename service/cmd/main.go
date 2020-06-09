package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"demo/pkg/authorizer"

	"github.com/julienschmidt/httprouter"
)

func main() {

	content, err := ioutil.ReadFile("../config/opa/data.json")
	data := string(content)
	log.Printf("%v", data)

	auth, err := authorizer.New("../config/opa", data)
	if err != nil {
		log.Printf("%v", err)
	}

	handler := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		allowed, err := auth.EvalRequest(r)
		if err != nil {
			http.Error(w, "Internal error failed to evaluate policy", http.StatusInternalServerError)
			return
		}
		if allowed {
			path := strings.Split(r.URL.Path, "/")[1:]
			w.Write([]byte("Hello Account: " + path[1]))
		} else {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
		}
	}

	router := httprouter.New()
	router.GET("/accounts/:id", handler)
	err = http.ListenAndServe("localhost:7777", router)
	if err != nil {
		log.Printf("%v", err)
	}
}
