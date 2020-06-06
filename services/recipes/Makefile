.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/process-recipes functions/dynamodb/recipeAdder/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/get-ingredients functions/api/ingredients/get/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/get-recipes functions/api/recipes/get/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/get-shopping-list functions/api/shopping-list/get/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
