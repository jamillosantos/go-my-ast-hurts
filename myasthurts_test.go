package myasthurts

import (
	"testing"

	"github.com/crowley-io/macchiato"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func TestMyASTHurts(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	macchiato.RunSpecs(t, "My AST Hurts Test Suite")
}
