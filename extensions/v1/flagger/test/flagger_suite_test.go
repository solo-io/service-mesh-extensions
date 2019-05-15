package test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	T *testing.T
)

func TestFlagger(t *testing.T) {
	T = t
	RegisterFailHandler(Fail)
	RunSpecs(t, "Flagger Suite")
}
