package memprovider_test

import (
	"context"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/open-feature/go-sdk/openfeature/memprovider"
	oft "openfeature-tests"
	"testing"
)

// Tests in this package are expected to fail when executed in parallel as they rely on the global state
//
//Test method: ensure a happens-before relationship between goroutines using latches
// to guarantee the following execution order:
//
// Test-one: setup Flag-A
// Test-two: setup Flag-B
// Test-two: execute code dependent of the Flag-B
// Test-one: execute code dependent of the Flag-A

var (
	latchA = oft.NewLatch()
	latchB = oft.NewLatch()
)

func TestOne_Memprovider(t *testing.T) {
	t.Parallel()

	testFlags := map[string]memprovider.InMemoryFlag{
		oft.FlagServiceAEnabled: oft.FlagA,
	}

	provider := memprovider.NewInMemoryProvider(testFlags)

	err := openfeature.SetProviderAndWait(provider)
	if err != nil {
		t.Fatalf("Failed to set the FlagA: %v", err)
	}

	latchB.CountDown()
	latchA.Wait()

	serviceA := oft.NewService(oft.FlagServiceAEnabled)

	if err = serviceA.Serve(context.Background()); err != nil {
		t.Fatalf("Error: %v", err)
	}
}

func TestTwo_Memprovider(t *testing.T) {
	t.Parallel()

	latchB.Wait()
	testFlags := map[string]memprovider.InMemoryFlag{
		oft.FlagServiceBEnabled: oft.FlagB,
	}

	provider := memprovider.NewInMemoryProvider(testFlags)

	if err := openfeature.SetProviderAndWait(provider); err != nil {
		t.Fatalf("Failed to set the FlagB: %v", err)
	}

	latchA.CountDown()
	serviceB := oft.NewService(oft.FlagServiceBEnabled)

	if err := serviceB.Serve(context.Background()); err != nil {
		t.Fatalf("Error: %v", err)
	}
}
