setup:
	go get github.com/dvyukov/go-fuzz/go-fuzz
	go get github.com/dvyukov/go-fuzz/go-fuzz-build

build:
	go-fuzz-build github.com/pascallouisperez/reflext/fuzz

run:
	mkdir -p examples/corpus
	cp initial-corpus/* examples/corpus/
	go-fuzz -bin=./fuzz-fuzz.zip -workdir=examples

clean:
	rm -f fuzz-fuzz.zip
	rm -Rf examples
