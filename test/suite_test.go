//go:build windows
// +build windows

package e2e

import (
    "testing"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

func TestTest(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Test Suite")
}
