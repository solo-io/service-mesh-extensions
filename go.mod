module github.com/solo-io/service-mesh-hub

go 1.13

require (
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/gogo/protobuf v1.3.1
	github.com/golang/mock v1.4.0
	github.com/golang/protobuf v1.3.4
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/helm/helm v2.13.1+incompatible
	github.com/huandu/xstrings v1.3.0 // indirect
	github.com/k0kubun/pp v3.0.1+incompatible // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/microcosm-cc/bluemonday v1.0.2
	github.com/mitchellh/reflectwalk v1.0.1 // indirect
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/pkg/errors v0.9.1
	github.com/rotisserie/eris v0.3.0
	github.com/russross/blackfriday v1.5.2
	github.com/solo-io/anyvendor v0.0.1
	github.com/solo-io/go-utils v0.14.1
	github.com/solo-io/solo-kit v0.13.2
	github.com/spf13/afero v1.2.2
	github.com/spf13/cobra v0.0.5
	go.uber.org/zap v1.13.0
	golang.org/x/crypto v0.0.0-20200221231518-2aa609cf4a9d // indirect
	golang.org/x/net v0.0.0-20200226121028-0de0cce0169b // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae // indirect
	golang.org/x/tools v0.0.0-20200226205201-eb7c56241bdb // indirect
	gopkg.in/AlecAivazis/survey.v1 v1.8.2
	gopkg.in/yaml.v2 v2.2.8 // indirect
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.3
	k8s.io/client-go v11.0.0+incompatible // indirect
	sigs.k8s.io/yaml v1.1.0
)

replace (
	// github.com/Azure/go-autorest/autorest has different versions for the Go
	// modules than it does for releases on the repository. Note the correct
	// version when updating.
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.0.0+incompatible
	github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.4.2
	github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309

	k8s.io/api => k8s.io/api v0.17.2
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.2
	k8s.io/apiserver => k8s.io/apiserver v0.17.2
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.17.2
	k8s.io/client-go => k8s.io/client-go v0.17.2
)
