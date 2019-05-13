.PHONY: generated-code
generated-code:
	go generate ./...
	protoc --gogo_out=. -I$(GOPATH)/src -I$(GOPATH)/src/github.com/gogo/protobuf -I$(GOPATH)/src/github.com/solo-io/service-mesh-hub api/v1/registry.proto

.PHONY: update-deps
update-deps:
	go get -u github.com/golang/mock/gomock
	go install github.com/golang/mock/mockgen
	go get -u github.com/gogo/protobuf/gogoproto
	go get -u github.com/gogo/protobuf/protoc-gen-gogo
	go get -u github.com/solo-io/solo-kit