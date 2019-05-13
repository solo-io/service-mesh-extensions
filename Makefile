.PHONY: generated-code
generated-code:
	protoc --gogo_out=. -I$(GOPATH)/src -I$(GOPATH)/src/github.com/gogo/protobuf -I$(GOPATH)/src/github.com/solo-io/service-mesh-hub api/v1/registry.proto

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
