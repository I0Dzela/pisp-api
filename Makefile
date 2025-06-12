clean:
	rm -rf gen

generate:
	mkdir -p gen/common gen/facekit gen/file gen/relation gen/user && go generate ./...

buildX:
	#cd src; go build -ldflags="-X 'build.Version=v1.0.0' -X 'build.User=$(id -u -n)' -X 'build.Time=$(date)'" -x -o ../bin/pisp main.go
	go build -ldflags="-X 'github.com/I0Dzela/pisp-specs/build.Version=3.0.2' -X 'github.com/I0Dzela/pisp-specs/build.User=$(shell whoami)' -X 'github.com/I0Dzela/pisp-specs/build.Time=$(shell TZ="CET" date +'%d.%m.%Y %H:%M:%S')'" -x -o ./bin/pisp-specs main.go

all: clean generate buildX
