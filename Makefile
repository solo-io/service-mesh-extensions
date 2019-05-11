.PHONY: init
init:
	go get -u github.com/solo-io/solo-kit

.PHONY: update-deps
update-deps:
	go get -u github.com/gogo/protobuf/gogoproto
	go get -u github.com/gogo/protobuf/protoc-gen-gogo

.PHONY: generated-code
generated-code:
	go run ci/pin_repos.go
	protoc --gogo_out=. -I$(GOPATH)/src -I$(GOPATH)/src/github.com/gogo/protobuf -I$(GOPATH)/src/github.com/solo-io/service-mesh-hub api/v1/registry.proto