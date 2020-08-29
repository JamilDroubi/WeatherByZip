.PHONY: all prog clean

all: prog clean

prog:
	go build weatherbyzip.go

clean:
	rm -rf *.core *.out *.o
