package repository

import (
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "bufio"
    "encoding/csv"
    "io"
    "os"
    . "companyProject/models"
    "strings"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "errors"
)

const (
	URI = "mongodb+srv://admin:admin@cluster0-nrksu.mongodb.net/test?retryWrites=true&w=majority"
	DATABASE = "companyDB"
    COLLECTION = "companies"
    FILEPATH = "C:/GoWork/src/companyProject/files/"
)

var db *mongo.Database
var collection *mongo.Collection

func Connect(collectionName string) error{
    if(collectionName == ""){
        collectionName = COLLECTION;
    }
    
	clientOptions := options.Client().ApplyURI(URI)

	client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        return err;
	}

	err = client.Ping(context.TODO(), nil)
    if err != nil {
        return err;
    }
    
    db = client.Database(DATABASE)
    collection = db.Collection(collectionName)

    return err;
}

func DropCollection() {
    collection.Drop(context.TODO());
}

func Create(company Company) (Company, error){
    res, err := collection.InsertOne(context.TODO(), company);
    if (err != nil) {
        return company, err;
    }

    id := res.InsertedID.(primitive.ObjectID);

    return ReadById(id);
}

func Read(filter bson.M) ([]Company, error){
    var companies []Company;
    var company Company;

    cur, err := collection.Find(context.Background(), filter)
    if err != nil {
        return companies, err;
    }

    defer cur.Close(context.Background())
    for cur.Next(context.Background()) {
    
        err := cur.Decode(&company)
        if err != nil { 
            return companies, err;
        }

        companies = append(companies, company);
    }

    return companies, err;
}

func ReadById(id primitive.ObjectID) (Company, error){
    var company Company;
    filter := bson.M{"_id": id}

    found := collection.FindOne(nil, filter);
    err := found.Decode(&company);
    if err != nil {
        return company, err;
    }

    return company, err;
}

func Update(companyFilter bson.M, update bson.M) error{
    _, err := collection.UpdateOne(context.TODO(), companyFilter, update)
    if err != nil {
        return err;
    }

    return err;
}

func Delete(companyFilter bson.M) error{
    _, err := collection.DeleteOne(context.TODO(), companyFilter)
    if err != nil {
        return err
    }
    
    return err;
}

func BuildFilterForAPI(name string, zip string) (bson.M){
    var filter  bson.M    

    if (name != "" && zip != "") {
        filter = bson.M{"company name":  bson.M{"$regex": ".*" + name + ".*"}, "zip code": zip}
    } else if (name != "") {
        filter = bson.M{"company name": bson.M{"$regex": ".*" + name + ".*"}}
    } else if (zip != "") {
        filter = bson.M{"zip code": zip}
    } else {
        filter = bson.M{}
    }

    return filter;	
}

func ValidateHeader(header []string) error{
    var err error

    if (len(header) != 3) {
        err = errors.New("Invalid number of headers");
    } else if (strings.ToLower(header[0]) != "name") {
        err = errors.New("The first header parameter must be name");
    } else if (strings.ToLower(header[1]) != "zip") {
        err = errors.New("The second header parameter must be zip");
    } else if (strings.ToLower(header[2]) != "website") {
        err = errors.New("The third header parameter must be website");
    }

    return err;
}

func ValidateContent(content []string) error{
    var err error

    companyName := content[0];
    companyZip := content[1];
    companyWebsite := content[2];

    if (companyName != strings.ToUpper(companyName)) {
        err = errors.New("The company name '" + companyName + "' must be in upper case");
    } else if (len(companyZip) != 5) {
        err = errors.New("The company zip '" + companyZip + "' must have 5 digits");
    } else if (companyWebsite != strings.ToLower(companyWebsite)) {
        err = errors.New("The company website '" + companyWebsite + "' must be in lower case");
    }

    return err;
}

func ValidateFileName(fileName string) error{
    var err error
    if (!strings.Contains(fileName, ".csv")){
        err = errors.New("Invalid file. The file must be a .csv");
    }

    return err;
}

func MergeFile(fileName string) error {
    err := ValidateFileName(fileName);
    if (err != nil) {
        return err;
    }

    fullFilePath := GetFullFilePath(fileName);
    csvFile, _ := os.Open(fullFilePath)
    reader := csv.NewReader(bufio.NewReader(csvFile))
    header, err := reader.Read()
    if (err != nil) {
        return err;
    }

    err = ValidateHeader(header);
    if (err != nil) {
        return err;
    }  

    for {
        line, err := reader.Read()
        if err == io.EOF {
            break
        } else if err != nil {
            return err;
        }      
        
        err = ValidateContent(line);
        if (err != nil) {
            continue;
        } 

        companyName := strings.ToLower(line[0]);    
        companyZip := line[1];    
        companyWebsite := strings.ToLower(line[2]);

        filter := BuildFilter(companyName, companyZip);
        update := bson.M{"$set": bson.M{"website": &companyWebsite}}
        
        Update(filter, update);
    }    

    return err;
}

func UploadFile(fileName string) error{
    fullFilePath := GetFullFilePath(fileName);
    csvFile, _ := os.Open(fullFilePath);
    reader := csv.NewReader(bufio.NewReader(csvFile));  
    _, err := reader.Read()
    if err != nil {
        return err;
    }

    var companies []Company;
    var company Company;

    for {
        line, err := reader.Read()
        if err == io.EOF {
            break
        } else if (err != nil) {
            return err;
        }

        company.Name = line[0];
        company.Zip = line[1];

        Create(company);
        
        companies = append(companies, company)
    }

    return err;
}

func GetFullFilePath(fileName string) string{
    return FILEPATH + fileName;
}

func BuildFilter(companyName string, companyZip string) bson.M {
    return bson.M{"company name": companyName, "zip code": companyZip};
}