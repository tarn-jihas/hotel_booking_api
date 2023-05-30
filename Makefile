build:
	@go build -o bin/api ;
run: build
	@./bin/api ;

seed:
	@go run scripts/seed.go ;

test: 
	@go test ./testing/handlers/... -v -count=1 ;
docker:
	echo "building docker file"
	@docker build -t api .
	echo "running API inside Docker container"
	@docker run -p 3333:3333 api