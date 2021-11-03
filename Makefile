lint:
	golint lib && golint lib/models && golint lib/internal/utils && golint lib/internal/messages && golint lib/internal/constants && golint examples

testLib:
	cd lib/ && go test -coverprofile=coverage.out

testLibModels:
	cd lib/models && go test -coverprofile=coverage.out 
	
testLibUtils:
	cd lib/internal/utils && go test -coverprofile=coverage.out 

test:
	make testLib
	make testLibModels
	make testLibUtils
	go test --coverprofile=coverage.out ./... && go tool cover -func=coverage.out
