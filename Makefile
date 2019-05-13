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

#----------------------------------------------------------------------------------
# Docs
#----------------------------------------------------------------------------------

site:
	if [ ! -d docs/themes ]; then git clone https://github.com/matcornic/hugo-theme-learn.git docs/themes/hugo-theme-learn; fi
	cd docs; hugo --config docs.toml

.PHONY: deploy-site
deploy-site: site
	firebase deploy --only hosting:mesh-market-place-docs

.PHONY: serve-site
serve-site: site
	cd docs; hugo --config docs.toml server -D
