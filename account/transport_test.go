package account

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log/level"
)

func TestDecodeCreateAccountRequestSuccess(t *testing.T) {
	var logger log.Logger
	var ctx context.Context
	handler := httptransport.NewServer(
		func(context.Context, interface{}) (interface{}, error) { return struct{}{}, nil },
		func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, errors.New("dang") },
		func(context.Context, http.ResponseWriter, interface{}) error { return nil },
	)
	server := httptest.NewServer(handler)
	defer server.Close()

	bodyRequest := createAccountRequest{
		ID:       "andi999",
		Balance:  100,
		Currency: "IDR",
	}
	byteArray, err := json.Marshal(bodyRequest)
	if err != nil {
		level.Error(logger).Log("error marshaling body request", err)
		os.Exit(-1)
	}

	r := httptest.NewRequest(http.MethodPost, server.URL+"/account/v1/account", bytes.NewBuffer(byteArray))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	dec, err := decodeCreateAccountRequest(ctx, r)
	if err != nil {
		level.Error(logger).Log("error decoding", err)
	}

	if want, have := fmt.Sprintf("%v", bodyRequest), fmt.Sprintf("%v", dec); want != have {
		t.Errorf("want %v, have %v", want, have)
	}
}

func TestDecodeCreateAccountRequestError(t *testing.T) {
	var logger log.Logger
	var ctx context.Context
	handler := httptransport.NewServer(
		func(context.Context, interface{}) (interface{}, error) { return struct{}{}, nil },
		func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, errors.New("dang") },
		func(context.Context, http.ResponseWriter, interface{}) error { return nil },
	)
	server := httptest.NewServer(handler)
	defer server.Close()

	bodyRequest := Account{
		ID:       "andi999",
		Balance:  100,
		Currency: "IDR",
	}
	byteArray, err := json.Marshal(bodyRequest)
	if err != nil {
		level.Error(logger).Log("error marshaling body request", err)
		os.Exit(-1)
	}

	r := httptest.NewRequest(http.MethodPost, server.URL+"/account/v1/account", bytes.NewBuffer(byteArray))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	dec, err := decodeCreateAccountRequest(ctx, r)
	if err != nil {
		level.Error(logger).Log("error decoding", err)
	}

	if want, have := fmt.Sprintf("%v", bodyRequest), fmt.Sprintf("%v", dec); want == have {
		t.Errorf("want %v, have %v", want, have)
	}
}
