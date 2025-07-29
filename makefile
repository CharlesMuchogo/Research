run:
	go run main.go

build:
	docker build -t research .

tag:
	docker tag research:latest charlesmuchogo/research:latest

push:
	docker push charlesmuchogo/research:latest

dev:
	@echo "starting dev..."
	@air