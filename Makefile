TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=controlplane.com
NAMESPACE=com
NAME=cpln
BINARY=terraform-provider-${NAME}
VERSION=1.2.14
OS_ARCH=linux_amd64

default: install

build:
	go build -o ${BINARY}

release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	# cd ./bin && tar -cvzf ${BINARY}_${VERSION}_darwin_amd64.tgz ${BINARY}_${VERSION}_darwin_amd64

	GOOS=darwin GOARCH=arm64 go build -o ./bin/${BINARY}_${VERSION}_darwin_arm64
	# cd ./bin && tar -cvzf ${BINARY}_${VERSION}_darwin_arm64.tgz ${BINARY}_${VERSION}_darwin_arm64

	# GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	# GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	# GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	# GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386

	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	# cd ./bin && tar -cvzf ${BINARY}_${VERSION}_linux_amd64.tgz ${BINARY}_${VERSION}_linux_amd64

	# GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	# GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	# GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	# GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	# GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386.exe

	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64.exe
	# cd ./bin && tar -cvzf ${BINARY}_${VERSION}_windows_amd64.zip ${BINARY}_${VERSION}_windows_amd64.exe

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	# cp ${BINARY} ~/tmp/terraform-cpln
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test: 
	go test -count=1 -i -v $(TEST) || exit 1                                                    
	echo $(TEST) | xargs -t -n4 go test -count=1 $(TESTARGS) -timeout=30s -parallel=4 -v                    

testacc: 
	# TF_ACC=1 CPLN_ORG=terraform-test-org CPLN_ENDPOINT=https://api.test.cpln.io CPLN_PROFILE=default CPLN_TOKEN= go test $(TEST) -v $(TESTARGS) -timeout 120m
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

clean:
	rm -rf ~/.terraform.d/plugins/${HOSTNAME}