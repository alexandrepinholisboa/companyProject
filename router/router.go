package router

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"companyProject/config/db"
)

func createJsonResponseError(w http.ResponseWriter, code int, msg string) {
	createJsonResponse(w, code, map[string]string{"error": msg})
}

func createJsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r);
	err := repository.UploadFile(params["id"]);
	if (err != nil) {
		createJsonResponseError(w, http.StatusInternalServerError, err.Error());
	} else {
		createJsonResponse(w, http.StatusOK, nil);
	}
}

func MergeFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r);
	err := repository.MergeFile(params["id"]);
	if (err != nil) {
		createJsonResponseError(w, http.StatusInternalServerError, err.Error());
	} else {
		createJsonResponse(w, http.StatusOK, nil);
	}

}

func Get(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name");
	zip := r.URL.Query().Get("zip");
	filter := repository.BuildFilterForAPI(name, zip);

	response, err := repository.Read(filter);
	if (err != nil) {
		createJsonResponseError(w, http.StatusInternalServerError, err.Error());
	} else if (len(response) == 0) {
		createJsonResponseError(w, http.StatusNotFound, "No Company was found on DB");
	} else {
		createJsonResponse(w, http.StatusOK, response);
	}
}