.PHONY: all get build docker deploy

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/iamrotator cmd/iamrotator/*.go
	chmod +x bin/iamrotator

get:
	go get -v -d ./...

docker:
	docker build --platform linux/arm64 -t loqutus/iamrotator:latest -f build/package/Dockerfile .
	docker push loqutus/iamrotator:latest

deploy:
	ssh ec2-user@iamrotator.govno.cloud 'docker stop iamrotator; docker rm iamrotator; docker pull loqutus/iamrotator:latest; docker run -d --name iamrotator -p 8080:8080 -v ~/config.yaml:/app/config.yaml loqutus/iamrotator:latest; docker logs iamrotator'

default: docker deploy