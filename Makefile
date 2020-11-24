# vim:set noexpandtab :
all:
	go build mopm.go
test:
	go test
dbuild:
	docker build -t mopm-test:latest .
drun:
	docker run -it mopm-test:latest /bin/ash
clean:
	rm -rf mopm
.PHONY: test clean dbuild drun

