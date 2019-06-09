package companyprojecttest

import (
	"testing"
	. "companyProject/models"
	"companyProject/config/db"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"github.com/gorilla/mux"	
	"companyProject/router"
	"net/http/httptest"
)

var companyRequest Company
var companyCreated Company

const (
	updatedWebsite = "updatedTest.com.br";
	companyName = "companyTest";
	companyZip = "12345";
	companyWebsite = "test.com.br";
)

func init(){
	companyRequest.Name = companyName;
	companyRequest.Zip = companyZip;
	companyRequest.Website = companyWebsite;
}

func TestConnection(t *testing.T) {
	repository.Connect("companiesTest");
}

func TestDropCollection(t *testing.T) {
	repository.DropCollection();
}

func TestCreate(t *testing.T) {	
	companyResponse, err := repository.Create(companyRequest);
	if (err != nil) {
		t.Errorf("An error occured within Create API: %v .", err.Error());
	}

	companyCreated = companyResponse;
}

func TestRead(t *testing.T) {	
	companyArray, err := repository.Read(bson.M{});
	if (err != nil) {
		t.Errorf("An error occured within Read API: %v .", err.Error());
	} 
	if(len(companyArray) == 0){
		t.Errorf("No record found on DB");
	}
}

func TestUpdate(t *testing.T) {
	filter := repository.BuildFilter(companyName, companyZip);
	update := bson.M{"$set": bson.M{"website": updatedWebsite}}
	
	err := repository.Update(filter, update); 	
	if (err != nil) {
		t.Errorf("An error occured within Update API: %v .", err.Error());
	} 
	
	companyResponse, err := repository.Read(filter);
	if (err != nil) {
		t.Errorf("An error occured within Read API: %v .", err.Error());
	} 

	company := companyResponse[0];
	if (company.Website != updatedWebsite) {
		t.Errorf("The Company website was not updated correctly");
	}
}

func TestDelete(t *testing.T) {
	filter := bson.M{"_id": companyCreated.ID}	
	err := repository.Delete(filter);
	if (err != nil) {
		t.Errorf("An error occured within Delete API: %v .", err.Error());
	}

	companyResponse, err := repository.Read(filter);
	if (err != nil) {
		t.Errorf("An error occured within Read API: %v .", err.Error());
	} else if(len(companyResponse) != 0) {
		t.Errorf("The Company was not deleted successfully");
	}
}

func TestUploadFile(t *testing.T) {
	repository.UploadFile("q1_catalog.csv");

	companyArray, err := repository.Read(bson.M{}); 
	if (err != nil) {
		t.Errorf("An error occured within Read API: %v .", err.Error());
	}

	companyArrayLength := len(companyArray);
	if(companyArrayLength != 44){
		t.Errorf("The UploadFile added %v records, expected: 44 .", companyArrayLength);
	}
}

func TestMergeValidFile(t *testing.T) {
	repository.MergeFile("q3_clientData.csv");
	
	companyArray, err := repository.Read(bson.M{"website": ""});
	if (err != nil) {
		t.Errorf("An error occured within Read API: %v .", err.Error());
	} 

	companyArrayLength := len(companyArray);
	if(companyArrayLength != 32){
		t.Errorf("Expected 32 records, found %v on DB", companyArrayLength);
	}
}

func TestRouterUploadFile(t *testing.T) {
	url := "/upload/q1_catalog.csv";

	request, err := http.NewRequest("POST", url, nil);
	if (err != nil) {
		t.Errorf("An error occured creating the request: %v .", err.Error());
	}
	
	responseRecorder := httptest.NewRecorder()

	muxRouter := mux.NewRouter();
	muxRouter.HandleFunc("/upload/{id}", router.UploadFile).Methods("POST");
	muxRouter.ServeHTTP(responseRecorder, request);

	if (responseRecorder.Code != 200){		
		t.Errorf("Response code %v, expected %v .", responseRecorder.Code, "200");
	}
}

func TestRouterMergeInvalidFile(t *testing.T) {
	url := "/merge/q2_clientData.csv";

	request, err := http.NewRequest("POST", url, nil);
	if (err != nil) {
		t.Errorf("An error occured creating the request: %v .", err.Error());
	}
	
	responseRecorder := httptest.NewRecorder()

	muxRouter := mux.NewRouter();
	muxRouter.HandleFunc("/merge/{id}", router.MergeFile).Methods("POST");
	muxRouter.ServeHTTP(responseRecorder, request);

	if (responseRecorder.Code != 500){		
		t.Errorf("Response code %v, expected 500 .", responseRecorder.Code);
	}

	companyArray, err := repository.Read(bson.M{"website": ""});
	if (err != nil) {
		t.Errorf("An error occured within Read API: %v .", err.Error());
	} 

	companyArrayLength := len(companyArray);
	if(companyArrayLength != 76){
		t.Errorf("Expected 76 records, found %v on DB", companyArrayLength);
	}
}

func TestRouterMergeTxtFile(t *testing.T) {
	url := "/merge/q4_clientData.txt";

	request, err := http.NewRequest("POST", url, nil);
	if (err != nil) {
		t.Errorf("An error occured creating the request: %v .", err.Error());
	}
	
	responseRecorder := httptest.NewRecorder()

	muxRouter := mux.NewRouter();
	muxRouter.HandleFunc("/merge/{id}", router.MergeFile).Methods("POST");
	muxRouter.ServeHTTP(responseRecorder, request);

	if (responseRecorder.Code != 500){		
		t.Errorf("Response code %v, expected 500 .", responseRecorder.Code);
	}
}

func TestRouterGet(t *testing.T) {
	url := "/company";

	request, err := http.NewRequest("GET", url, nil);
	if (err != nil) {
		t.Errorf("An error occured creating the request: %v .", err.Error());
	}

	responseRecorder := httptest.NewRecorder()
	http.HandlerFunc(router.Get).ServeHTTP(responseRecorder, request)

	if (responseRecorder.Code != 200){		
		t.Errorf("Response code %v, expected %v .", responseRecorder.Code, "200");
	}

	companyArray, err := repository.Read(bson.M{});
	if (err != nil) {
		t.Errorf("An error occured within Read API: %v .", err.Error());
	} 

	companyArrayLength := len(companyArray);
	if(companyArrayLength != 88){
		t.Errorf("Expected 88 records, found %v on DB", companyArrayLength);
	}
}

func TestRouterGetByName(t *testing.T) {
	url := "/company/?name=sales";

	request, err := http.NewRequest("GET", url, nil);
	if (err != nil) {
		t.Errorf("An error occured creating the request: %v .", err.Error());
	}
	
	responseRecorder := httptest.NewRecorder()

	muxRouter := mux.NewRouter();
	muxRouter.HandleFunc("/company/", router.Get).Methods("GET");
	muxRouter.ServeHTTP(responseRecorder, request);

	if (responseRecorder.Code != 200){		
		t.Errorf("Response code %v, expected %v .", responseRecorder.Code, "200");
	}	
}

func TestRouterGetByZip(t *testing.T) {
	url := "/company/?zip=57761";

	request, err := http.NewRequest("GET", url, nil);
	if (err != nil) {
		t.Errorf("An error occured creating the request: %v .", err.Error());
	}
	
	responseRecorder := httptest.NewRecorder()

	muxRouter := mux.NewRouter();
	muxRouter.HandleFunc("/company/", router.Get).Methods("GET");
	muxRouter.ServeHTTP(responseRecorder, request);

	if (responseRecorder.Code != 200){		
		t.Errorf("Response code %v, expected %v .", responseRecorder.Code, "200");
	}	
}

func TestRouterGetByNameAndZip(t *testing.T) {
	url := "/company/?name=post&zip=57761";

	request, err := http.NewRequest("GET", url, nil);
	if (err != nil) {
		t.Errorf("An error occured creating the request: %v .", err.Error());
	}
	
	responseRecorder := httptest.NewRecorder()

	muxRouter := mux.NewRouter();
	muxRouter.HandleFunc("/company/", router.Get).Methods("GET");
	muxRouter.ServeHTTP(responseRecorder, request);

	if (responseRecorder.Code != 200){		
		t.Errorf("Response code %v, expected %v .", responseRecorder.Code, "200");
	}	
}