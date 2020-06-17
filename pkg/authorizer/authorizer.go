package authorizer

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/storage/inmem"
	"github.com/open-policy-agent/opa/util"
)

type opaInput struct {
	Method   interface{} `json:"method"`
	Path     interface{} `json:"path"`
	Roles    interface{} `json:"roles"`
	Region   interface{} `json:"region"`
	UserName interface{} `json:"userName"`
}

type Accounts struct {
	Accounts []Account `json:"accounts"`
}

type Account struct {
	Id       int    `json:"ID"`
	Username string `json:"username"`
	Region   string `json:"region"`
}

type OPAAuthorizer struct {
	store storage.Store
	txn   storage.Transaction
	model func(r *rego.Rego)
}

// Returns a new OPA object
func New(opaDirectory string, dataStoreBytes []byte) (opaAuthorizer OPAAuthorizer, err error) {
	ctx := context.Background()
	var dataStore map[string]interface{}

	err = util.UnmarshalJSON(dataStoreBytes, &dataStore)
	if err != nil {
		fmt.Println("Failed to parse data store json")
	}

	store := inmem.NewFromObject(dataStore)
	txn, err := store.NewTransaction(ctx, storage.WriteParams)
	if err != nil {
		return opaAuthorizer, err
	}
	opaAuthorizer = OPAAuthorizer{
		txn:   txn,
		store: store,
		model: rego.Load([]string{opaDirectory}, nil),
	}
	return opaAuthorizer, nil
}

func (opaAuth *OPAAuthorizer) EvalRequest(req *http.Request) (allowed bool, err error) {
	ctx := context.Background()
	r := rego.New(
		rego.Query("allowed = data.demo.allow"),
		opaAuth.model,
		rego.Store(opaAuth.store),
		rego.Transaction(opaAuth.txn),
	)
	query, err := r.PrepareForEval(ctx)
	if err != nil {
		return false, err
	}
	input := opaAuth.convertRequestToInput(req)
	rs, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return false, err
	}
	allowed, ok := rs[0].Bindings["allowed"].(bool)
	if !ok {
		return false, errors.New("Failed to convert allowed to bool")
	}
	return allowed, err
}

func (opaAuth *OPAAuthorizer) convertRequestToInput(req *http.Request) opaInput {
	input := opaInput{
		Method:   req.Method,
		Path:     strings.Split(req.URL.Path, "/")[1:],
		Region:   req.Header.Get("region"),
		Roles:    strings.Split(req.Header.Get("roles"), ","),
		UserName: req.Header.Get("username"),
	}
	return input
}
