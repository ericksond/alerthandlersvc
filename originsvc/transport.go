package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
)

func makeProcessAlertEndpoint(svc OriginService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(processalertRequest)
		v, err := svc.ProcessAlert(req.S)
		if err != nil {
			return processalertResponse{v, err.Error()}, nil
		}
		return processalertResponse{v, ""}, nil
	}
}

func makeListEndpoint(svc OriginService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listRequest)
		v, err := svc.List(req.S)
		if err != nil {
			return listResponse{v, err.Error()}, nil
		}
		return listResponse{v, ""}, nil
	}
}

func decodeProcessAlertRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request processalertRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request listRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

type processalertRequest struct {
	S string `json:"search_name"`
}

type processalertResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"`
}

type listRequest struct {
	S string `json:"s"`
}

type listResponse struct {
	Alerts map[string]interface{}
	Err    string `json:"err,omitempty"`
}
