.PHONY: generated-code
generated-code:
	go generate ./...
	protoc --gogo_out=. -I$(GOPATH)/src -I$(GOPATH)/src/github.com/gogo/protobuf -I$(GOPATH)/src/github.com/gogo/protobuf/protobuf -I$(GOPATH)/src/github.com/solo-io/service-mesh-hub api/v1/registry.proto

.PHONY: update-deps
update-deps:
	go get github.com/golang/mock/gomock
	go get github.com/golang/mock/mockgen # fix vendoring problem also surfaced here: https://github.com/openshift/openshift-azure/issues/1582
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
