package internal

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"database/sql"

	"strconv"

	transport "github.com/go-kit/kit/transport/http"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func decodeAuthenticateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req AuthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func decodeHomeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return HomeRequest{}, nil
}

func decodeGetUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	i, ok := vars["id"]
	if !ok {
		return nil, errors.New("bad route")
	}

	id, err := strconv.Atoi(i)
	if err != nil {
		return nil, errors.New("bad route")
	}

	return GetUserRequest{Id: uint(id)}, nil
}

func decodeGetUsersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return GetUsersRequest{}, nil
}

func decodeCreateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	e, ok := response.(errorer)
	if ok && e.error() != nil {
		// encode errors from business logic
		switch e.error() {
		case ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
		case ErrInvalidArgument:
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return json.NewEncoder(w).Encode(map[string]interface{}{
			"error": e.error().Error(),
		})
	}

	return json.NewEncoder(w).Encode(response)
}

func NewHTTPHandler() http.Handler {
	// TODO: Get from config
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3326)/gokit")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Connected to db")

	svc := NewAuthService(db)

	r := mux.NewRouter()
	r.Methods("GET").Path("/").Handler(transport.NewServer(MakeHomeEndpoint(), decodeHomeRequest, encodeResponse))
	r.Methods("POST").Path("/auth").Handler(transport.NewServer(MakeAuthenticateEndpoint(svc), decodeAuthenticateRequest, encodeResponse))
	r.Methods("GET").Path("/user/{id}").Handler(transport.NewServer(MakeGetUserEndpoint(svc), decodeGetUserRequest, encodeResponse))
	r.Methods("GET").Path("/users").Handler(transport.NewServer(MakeGetUsersEndpoint(svc), decodeGetUsersRequest, encodeResponse))
	r.Methods("POST").Path("/user").Handler(transport.NewServer(MakeCreateUserEndpoint(svc), decodeCreateUserRequest, encodeResponse))

	return r
}
