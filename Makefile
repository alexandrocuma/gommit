.PHONY:  format

install:
	go install .
	
format:
	go fmt ./...
