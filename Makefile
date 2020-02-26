#----------------------------------------------------------------------------------
# Base
#----------------------------------------------------------------------------------

ROOTDIR := $(shell pwd)
OUTPUT_DIR ?= $(ROOTDIR)/_output
VERSION ?= $(shell echo $(TAGGED_VERSION) | cut -c 2-)
LDFLAGS := "-X github.com/solo-io/service-mesh-hub/pkg/internal/version.Version=$(VERSION)"
GCFLAGS := all="-N -l"

.PHONY: generated-code
generated-code:
	protoc --gogo_out=Mgoogle/protobuf/timestamp.proto=github.com/golang/protobuf/ptypes/timestamp:. -I$(GOPATH)/src -I$(GOPATH)/src/github.com/gogo/protobuf -I$(GOPATH)/src/github.com/gogo/protobuf/protobuf -I$(GOPATH)/src/github.com/solo-io/service-mesh-hub api/v1/registry.proto
	go generate ./...
	gofmt -w $(shell ls -d -- */ | grep -v vendor) && goimports -w $(shell ls -d -- */ | grep -v vendor)

.PHONY: update-deps
update-deps:
	GO111MODULE=off go get golang.org/x/tools/cmd/goimports
	GO111MODULE=off go get github.com/golang/mock/gomock
	GO111MODULE=off go get github.com/golang/mock/mockgen # fix vendoring problem also surfaced here: https://github.com/openshift/openshift-azure/issues/1582
	GO111MODULE=off go install github.com/golang/mock/mockgen
	GO111MODULE=off go get -u github.com/gogo/protobuf/gogoproto
	GO111MODULE=off go get -u github.com/gogo/protobuf/protoc-gen-gogo

#----------------------------------------------------------------------------------
# hubctl
#----------------------------------------------------------------------------------
CLI_DIR=./pkg/cli

$(OUTPUT_DIR)/hubctl: $(SOURCES)
	go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(CLI_DIR)/cmd/main.go

$(OUTPUT_DIR)/hubctl-linux-amd64: $(SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(CLI_DIR)/cmd/main.go

$(OUTPUT_DIR)/hubctl-darwin-amd64: $(SOURCES)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ $(CLI_DIR)/cmd/main.go

.PHONY: hubctl
hubctl: $(OUTPUT_DIR)/hubctl
.PHONY: hubctl-linux-amd64
hubctl-linux-amd64: $(OUTPUT_DIR)/hubctl-linux-amd64
.PHONY: hubctl-darwin-amd64
hubctl-darwin-amd64: $(OUTPUT_DIR)/hubctl-darwin-amd64

.PHONY: build-cli
build-cli: hubctl-linux-amd64 hubctl-darwin-amd64

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

.PHONY: clean-site
clean:
	rm -fr ./site ./resources

# Uses https://github.com/gjtorikian/html-proofer
# Does not require running site; just make sure you generate the site and then run it
# Install with gem install html-proofer
# Another option we could use is wget: https://www.digitalocean.com/community/tutorials/how-to-find-broken-links-on-your-website-using-wget-on-debian-7
.PHONY: check-links
check-links:
	cd docs; hugo --config docs.toml check
	htmlproofer ./site/ --allow-hash-href --alt-ignore "/img/Gloo-01.png" --url-ignore "/localhost/,/github.com/solo-io/solo-projects/,/developers.google.com/,/getgrav.org/,/github.com/solo-io/gloo/projects/,/developer.mozilla.org/"
