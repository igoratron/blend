.PHONY: build clean deploy

build:
	cd github.com/igoratron/blend/; env GOOS=linux go build -ldflags="-s -w" -o ../../../bin/process-recipes recipeAdder/main.go
	cd github.com/igoratron/blend/; env GOOS=linux go build -ldflags="-s -w" -o ../../../bin/get-shopping-list getShoppingList/main.go
	cd github.com/igoratron/blend/; env GOOS=linux go build -ldflags="-s -w" -o ../../../bin/get-ingredients getIngredients/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
