build:
	@go build -o captcha.out server/main.go

start:
	@go run server/main.go --stype redis --saddress localhost:6379 --expire 300

.PHONY: start
