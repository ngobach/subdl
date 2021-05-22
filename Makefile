INSTALL_DST := /usr/local/bin

default: build

dist:
	@rm -rf dist
	@mkdir dist
	@echo "\"dist\" created"

dist/subdl: dist sub main.go
	@go build -o $@ .
	@strip $@
	@echo "Build $@ done"

dist/subdl.exe: dist sub main.go
	@GOOS=windows go build -o $@ .
	@strip $@
	@echo "Build $@ done"

build: dist/subdl dist/subdl.exe

install: dist/subdl
	cp dist/subdl $(INSTALL_DST)/subdl