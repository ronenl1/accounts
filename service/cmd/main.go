package main

import (
	"log"
	"net/http"

	"demo/pkg/authorizer"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()

	auth, err := authorizer.New("../config/opa", "")
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
			w.Write([]byte("kakoi"))
		} else {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
		}
	}
	router.GET("/accounts/:id", handler)
	err = http.ListenAndServe("localhost:7777", router)
	if err != nil {
		log.Printf("%v", err)
	}
}
