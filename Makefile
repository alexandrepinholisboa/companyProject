#Challenge Makefile

start:
go run companyProject/

check:
go test -v companyProject/tests

#setup:
go get -u github.com/gorilla/mux
go get go.mongodb.org/mongo-driver/bson
go get go.mongodb.org/mongo-driver/bson/primitive
go get go.mongodb.org/mongo-driver/mongo
go get go.mongodb.org/mongo-driver/mongo/options
go install companyProject
