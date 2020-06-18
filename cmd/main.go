package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"demo/pkg/authorizer"

	"github.com/julienschmidt/httprouter"
	"github.com/open-policy-agent/opa/util"
)

func main() {

	dataStoreBytes, _ := ioutil.ReadFile("./config/opa/data.json")

	var dataStore map[string]interface{}

	err := util.UnmarshalJSON(dataStoreBytes, &dataStore)
	if err != nil {
		return
	}

	auth, err := authorizer.New("./config/opa", dataStore)

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
			fmt.Println(path[1])
			someMap := dataStore["accounts"]
			fmt.Println(someMap)
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
