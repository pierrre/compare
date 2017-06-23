test:
	go test -v

lint:
	go get -v github.com/alecthomas/gometalinter
	gometalinter --install --update --no-vendored-linters
	gometalinter --enable-all -D dupl -D lll -D gas -D goconst -D gocyclo -D gotype -D interfacer -D safesql -D test -D testify -D vetshadow\
	 --tests --deadline=10m --concurrency=2

.PHONY: test lint
