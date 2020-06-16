package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"demo/pkg/authorizer"

	"github.com/julienschmidt/httprouter"
)

func main() {

	content, err := ioutil.ReadFile("./config/opa/data.json")

	auth, err := authorizer.New("./config/opa", content)
	if err != nil {
		log.Printf("%v", err)
	} else {
		fmt.Println("OPA engine is up!")
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
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	}

	router := httprouter.New()
	router.GET("/accounts/:id", handler)
	err = http.ListenAndServe("localhost:7777", router)
	if err != nil {
		log.Printf("%v", err)
	}
}
