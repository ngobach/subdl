.PHONY: dist

default: build

dist:
	@rm -rf dist
	@mkdir dist
	@echo "\"dist\" created"

dist/subdl: dist
	@go build -o $@ .
	@strip $@
	@echo "Build done"

build: dist/subdl
