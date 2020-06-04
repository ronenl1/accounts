package authorizer

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/storage/inmem"
)

type opaInput struct {
	Method  interface{} `json:"method"`
	Path    interface{} `json:"path"`
	Headers interface{} `json:"headers"`
}

type OPAAuthorizer struct {
	store storage.Store
	txn   storage.Transaction
	model func(r *rego.Rego)
}

func New(opaDirectory, storeData string) (opaAuthorizer OPAAuthorizer, err error) {
	ctx := context.Background()
	buffer := bytes.NewBufferString(storeData)
	store := inmem.NewFromReader(buffer)
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
		rego.Query("allowed = data.ast_jwt.allow"),
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
		Method:  req.Method,
		Path:    strings.Split(req.URL.Path, "/")[1:],
		Headers: req.Header,
	}
	return input
}
