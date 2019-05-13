package kustomize_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	T *testing.T
)

func TestKustomize(t *testing.T) {
	T = t
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kustomize Suite")
}
