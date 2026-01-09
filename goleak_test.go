package openfeature_examples

import (
	"fmt"
	"testing"

	"github.com/open-feature/go-sdk/openfeature/memprovider"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	// make sure we don't leak goroutines after tests in this package have
	// finished, which means we haven't leaked contexts either
	goleak.VerifyTestMain(m)
}

func TestGoleak(t *testing.T) {
	flag := memprovider.InMemoryFlag{
		Key: "DUMMY_KEY",
	}

	fmt.Print(flag.Key)
}
