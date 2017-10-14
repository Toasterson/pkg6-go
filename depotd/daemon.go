package main

import (
	"github.com/gorilla/mux"
	"github.com/minio/minio/pkg/http"
)

var router *mux.Router

func init() {
	router = mux.NewRouter()
	router.HandleFunc("/versions/0/", handleVersionsV0)
	router.HandleFunc("/{publisher}/search/0/{query}", handleSearchV0)
	router.HandleFunc("/{publisher}/search/1/{casesensitive}_{returntype}_{maxreturn}_{startreturn}_{query}", handleSearchV1)
	router.HandleFunc("/{publisher}/catalog/0/{catalog}", handleCatalogV0)
	router.HandleFunc("/{publisher}/catalog/1/{catalog}", handleCatalogV1)
}

func main() {
	server := http.NewServer([]string{"0.0.0.0:8080"}, router, nil)
	err := server.ListenAndServe()
	if err != nil{
		panic(err)
	}
}

