.PHONY: build clean deploy

build:
	cd github.com/igoratron/blend/; env GOOS=linux go build -ldflags="-s -w" -o ../../../bin/recipeAdder recipeAdder/main.go
	cd github.com/igoratron/blend/; env GOOS=linux go build -ldflags="-s -w" -o ../../../bin/shoppingListCreator getShoppingList/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
