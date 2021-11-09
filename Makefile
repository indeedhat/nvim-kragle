.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -o . ./...

.PHONY: build-osx
build-osx:
	CGO_ENABLED=0 GOOS=darwin go build -o . ./...

install:
	OS=`uname`
	if [ $OS == "Linux" ]; then
		make build
	elif [ $OS == "Darwin" ]; then
		make build-osx
	else
		echo "OS not supported"
		exit 1
	fi
