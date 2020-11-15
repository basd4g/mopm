all:
	go build mopm.go
test:
	go test
clean:
	rm -rf mopm
.PHONY: test clean

