run:
	go run . search Untouchable

build:
	rm -rf dist && mkdir dist
	go build -o ./dist/subdl .
	# strip ./dist/subdl
