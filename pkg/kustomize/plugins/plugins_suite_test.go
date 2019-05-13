package plugins

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	T *testing.T
)

func TestPlugins(t *testing.T) {
	T = t
	RegisterFailHandler(Fail)
	RunSpecs(t, "Plugins Suite")
}
