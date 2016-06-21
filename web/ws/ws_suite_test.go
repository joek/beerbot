package ws_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestWs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ws Suite")
}
