package test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	T *testing.T
)

func TestKiali(t *testing.T) {
	T = t
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kiali Suite")
}
