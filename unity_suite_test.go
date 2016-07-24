package unity_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestUnity(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Unity Suite")
}
