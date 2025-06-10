clean:
	rm -rf gen

generate:
	mkdir -p gen/api && go generate ./...

all: clean generate
