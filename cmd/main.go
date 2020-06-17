package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"demo/pkg/authorizer"

	"github.com/julienschmidt/httprouter"
)

func main() {

	dataStoreBytes, _ := ioutil.ReadFile("./config/opa/data.json")

	auth, err := authorizer.New("./config/opa", dataStoreBytes)

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
			_, err = w.Write([]byte("Hello " + r.Header.Get("username")))
			if err != nil {
				return
			}
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
