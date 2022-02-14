package payment

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

func TestDecodeCreatePaymentRequestSuccess(t *testing.T) {
	var logger log.Logger
	var ctx context.Context
	handler := httptransport.NewServer(
		func(context.Context, interface{}) (interface{}, error) { return struct{}{}, nil },
		func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, errors.New("dang") },
		func(context.Context, http.ResponseWriter, interface{}) error { return nil },
	)
	server := httptest.NewServer(handler)
	defer server.Close()

	bodyRequest := sendPaymentRequest{
		AccountID: "andi999",
		Amount:    20,
		ToAccount: "maria007",
	}
	byteArray, err := json.Marshal(bodyRequest)
	if err != nil {
		level.Error(logger).Log("error marshaling body request", err)
		os.Exit(-1)
	}

	r := httptest.NewRequest(http.MethodPost, server.URL+"/account/v1/account", bytes.NewBuffer(byteArray))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	dec, err := decodeSendPaymentRequest(ctx, r)
	if err != nil {
		level.Error(logger).Log("error decoding", err)
	}

	if want, have := fmt.Sprintf("%v", bodyRequest), fmt.Sprintf("%v", dec); want != have {
		t.Errorf("want %v, have %v", want, have)
	}
}

func TestDecodeCreatePaymentRequestError(t *testing.T) {
	var logger log.Logger
	var ctx context.Context
	handler := httptransport.NewServer(
		func(context.Context, interface{}) (interface{}, error) { return struct{}{}, nil },
		func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, errors.New("dang") },
		func(context.Context, http.ResponseWriter, interface{}) error { return nil },
	)
	server := httptest.NewServer(handler)
	defer server.Close()

	bodyRequest := Payment{
		AccountID: "andi999",
		Amount:    20,
		ToAccount: "maria007",
	}
	byteArray, err := json.Marshal(bodyRequest)
	if err != nil {
		level.Error(logger).Log("error marshaling body request", err)
		os.Exit(-1)
	}

	r := httptest.NewRequest(http.MethodPost, server.URL+"/account/v1/account", bytes.NewBuffer(byteArray))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	dec, err := decodeSendPaymentRequest(ctx, r)
	if err != nil {
		level.Error(logger).Log("error decoding", err)
	}

	if want, have := fmt.Sprintf("%v", bodyRequest), fmt.Sprintf("%v", dec); want == have {
		t.Errorf("want %v, have %v", want, have)
	}
}
