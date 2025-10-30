package package_a_test

import (
	"fmt"
	"github.com/open-feature/go-sdk/openfeature"
	oftesting "github.com/open-feature/go-sdk/openfeature/testing"
	oft "openfeature-tests"
	"os"
	"testing"
)

var (
	latchTestA = oft.NewLatch()
	latchTestB = oft.NewLatch()

	// Since openfeature relies on global state,
	// TestProvider should be a global instance shared between tests
	// Under the hood it utilises a custom _goroutine local_ storage to manage state on a per-test basis
	provider = oftesting.NewTestProvider()
)

func TestMain(m *testing.M) {
	fmt.Println("Setting up test environment...")

	if err := openfeature.SetProvider(provider); err != nil {
		panic(err)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}
