.PHONY: tidy verify update deps

tidy:
	go mod tidy

verify:
	go mod verify

update:
	go get -u ./...
	go mod tidy

deps: tidy verify
