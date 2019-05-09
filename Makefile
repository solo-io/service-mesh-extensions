.PHONY: generated-code
generated-code:
	protoc --gogo_out=generated -I$(GOPATH)/src -I$(GOPATH)/src/github.com/gogo/protobuf -I$(GOPATH)/src/github.com/solo-io/service-mesh-hub api/v1/registry.proto