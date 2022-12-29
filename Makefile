BINARY_NAME=poco

WINDOWS=latest/windows/${BINARY_NAME}.exe
DARWIN=latest/mac/${BINARY_NAME}
LINUX=latest/linux/${BINARY_NAME}

build:
	mkdir -p latest/{darwin,linux,windows}
	jGOARCH=amd64 GOOS=darwin go build -o ${DARWIN} main.go
	GOARCH=amd64 GOOS=linux go build -o ${LINUX} main.go
	GOARCH=amd64 GOOS=windows go build -o ${WINDOWS} main.go

run:
	./${BINARY_NAME}

build_and_run: build run

clean:
	go clean
	rm ${WINDOWS}
	rm ${LINUX}
	rm ${DARWIN}

