package beerbot_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBeerbot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Beerbot Suite")
}
