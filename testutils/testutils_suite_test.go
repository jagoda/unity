package testutils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTestutils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testutils Suite")
}
