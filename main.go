package main

import (
    "log"
    "github.com/gorilla/mux"
	"companyProject/config/db"
    "companyProject/router"
    "net/http"
)

func main() {
    repository.Connect("");
    
    muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/company", router.Get).Methods("GET")
	muxRouter.HandleFunc("/upload/{id}", router.UploadFile).Methods("POST")
	muxRouter.HandleFunc("/merge/{id}", router.MergeFile).Methods("POST")
    log.Fatal(http.ListenAndServe(":8000", muxRouter))

}