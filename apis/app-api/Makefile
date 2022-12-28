install-dependencies:
	npm ci

build:
	export GO111MODULE=on
	export CGO_ENABLED=0
	go mod download

	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/api/app/auth/signup handler/auth/signup/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/api/app/auth/signin handler/auth/signin/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/api/app/auth/authorizer handler/auth/authorizer/main.go

clean:
	rm -rf ./bin

deploy:	
	serverless create-cert --verbose
	serverless create_domain --verbose
	serverless deploy --verbose
