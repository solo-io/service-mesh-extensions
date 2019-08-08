package registry_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const DefaultRemoteUrl = "https://storage.googleapis.com/sm-marketplace-registry/extensions4.yaml"

func TestRegistry(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry Suite")
}
