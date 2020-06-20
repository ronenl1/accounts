package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
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
			accountsMap := dataStore["accounts"]
			selectedAccount := getAccountById(accountsMap, path[1])
			selectedAccountJson, _ := json.Marshal(selectedAccount)
			_, err = w.Write(selectedAccountJson)
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

func getAccountById(accountsMap interface{}, accountId string) interface{} {
	return reflect.ValueOf(accountsMap).MapIndex(reflect.ValueOf(accountId)).Elem().Interface()
}
