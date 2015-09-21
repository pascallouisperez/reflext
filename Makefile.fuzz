default: clean build run

setup:
	go get github.com/dvyukov/go-fuzz/go-fuzz
	go get github.com/dvyukov/go-fuzz/go-fuzz-build

build:
	cp fuzz.go.ignore fuzz.go
	go-fuzz-build github.com/pascallouisperez/reflext

run:
	mkdir -p examples/corpus
	cp initial-corpus/* examples/corpus/
	go-fuzz -bin=./reflext-fuzz.zip -workdir=examples

clean:
	rm -f fuzz.go reflext-fuzz.zip
	rm -Rf examples
