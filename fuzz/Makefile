setup:
	go get github.com/dvyukov/go-fuzz/go-fuzz
	go get github.com/dvyukov/go-fuzz/go-fuzz-build

build:
	go-fuzz-build github.com/pascallouisperez/reflext/fuzz

run:
	go-fuzz -bin=./fuzz-fuzz.zip -workdir=examples

clean:
	rm fuzz-fuzz.zip
	rm -Rf examples
	git checkout -- examples
