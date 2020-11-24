# vim:set noexpandtab :
all:
	cd src && go build -o ../mopm
test:
	cd src && go test
dbuild:
	docker build -t mopm-test:latest .
dtest:
	make dbuild
	docker run mopm-test:latest
drun:
	docker run -it mopm-test:latest /bin/ash
clean:
	rm -rf mopm
.PHONY: test clean dbuild drun

