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

type OPAAuthorizer struct {
	store storage.Store
	txn   storage.Transaction
	model func(r *rego.Rego)
}

func New(opaDirectory string, storeData []byte) (opaAuthorizer OPAAuthorizer, err error) {
	ctx := context.Background()
	fmt.Printf("%v", string(storeData))

	var json map[string]interface{}
	err = util.UnmarshalJSON(storeData, &json)
	if err != nil {
		// Handle error.
	}
	store := inmem.NewFromObject(json)
	// buffer := bytes.NewBufferString(storeData)
	// store := inmem.NewFromObject(buffer)
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
		return false, errors.New("failed to convert allowed to bool")
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
