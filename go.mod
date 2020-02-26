module github.com/solo-io/service-mesh-hub

go 1.13

require (
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/gogo/protobuf v1.3.1
	github.com/golang/mock v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/google/go-github v17.0.0+incompatible

	github.com/helm/helm v2.13.1+incompatible
	github.com/microcosm-cc/bluemonday v1.0.2
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/pkg/errors v0.8.1
	github.com/rotisserie/eris v0.1.1
	github.com/russross/blackfriday v1.5.2
	github.com/solo-io/go-utils v0.13.0
	github.com/solo-io/solo-kit v0.13.0
	github.com/spf13/afero v1.2.2
	github.com/spf13/cobra v0.0.5
	go.uber.org/zap v1.13.0
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	gopkg.in/AlecAivazis/survey.v1 v1.8.2
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
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
