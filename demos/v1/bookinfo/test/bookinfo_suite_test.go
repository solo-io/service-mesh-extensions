package test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	T *testing.T
)

func TestBookinfo(t *testing.T) {
	T = t
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bookinfo Suite")
}
