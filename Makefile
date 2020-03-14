# Go parameters
GOCMD=go
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

test: 
	$(GOTEST) -v ./...
test-watch:
	ls *.go | entr sh -c '$(GOTEST) -v ./... -coverprofile=cover.out'	
