build:
	go-fuzz-build github.com/pascallouisperez/reflext/fuzz

run:
	go-fuzz -bin=./fuzz-fuzz.zip -workdir=examples

clean:
	rm fuzz-fuzz.zip
	rm -Rf examples
	git checkout -- examples
