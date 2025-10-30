package package_b_test

import (
	"context"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/open-feature/go-sdk/openfeature/memprovider"
	oft "openfeature-tests"
	"testing"
)

// Tests in this package are expected to pass when executed in parallel
// as they rely on the TestProvider that utilises a custom _goroutine local_ storage
// to manage state on a per-test basis.
//
// Test method: ensure a happens-before relationship between goroutines using latches
// to guarantee the following execution order:
//
// Test-one: setup Flag-A
// Test-two: setup Flag-B
// Test-two: execute code dependent of the Flag-B
// Test-one: execute code dependent of the Flag-A

func TestOne_PackageA(t *testing.T) {
	t.Parallel()

	testFlags := map[string]memprovider.InMemoryFlag{
		oft.FlagServiceAEnabled: oft.FlagA,
	}

	provider.UsingFlags(t, testFlags)

	err := openfeature.SetProviderAndWait(provider)
	if err != nil {
		t.Fatalf("failed to set provider: %v", err)
	}

	latchTestB.CountDown()
	latchTestA.Wait()

	serviceA := oft.NewService(oft.FlagServiceAEnabled)

	if err = serviceA.Serve(context.Background()); err != nil {
		t.Fatalf("failed to serve: %v", err)
	}
}

func TestTwo_PackageB(t *testing.T) {
	t.Parallel()

	latchTestB.Wait()
	testFlags := map[string]memprovider.InMemoryFlag{
		oft.FlagServiceBEnabled: oft.FlagB,
	}

	provider.UsingFlags(t, testFlags)

	if err := openfeature.SetProviderAndWait(provider); err != nil {
		t.Fatalf("failed to set provider: %v", err)
	}

	latchTestA.CountDown()
	serviceB := oft.NewService(oft.FlagServiceBEnabled)

	if err := serviceB.Serve(context.Background()); err != nil {
		t.Fatalf("failed to serve: %v", err)
	}
}
