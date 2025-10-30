package tests

import (
	"context"
	"fmt"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/open-feature/go-sdk/openfeature/memprovider"
	"sync"
)

const (
	FlagServiceAEnabled = "service-a-enabled"
	FlagServiceBEnabled = "service-b-enabled"
	FlagServiceCEnabled = "service-c-enabled"
	FlagServiceDEnabled = "service-d-enabled"
)

var (
	FlagA = memprovider.InMemoryFlag{

		Key:            FlagServiceAEnabled,
		State:          memprovider.Enabled,
		DefaultVariant: "on",
		Variants: map[string]any{
			"on":  true,
			"off": false,
		}}

	FlagB = memprovider.InMemoryFlag{
		Key:            FlagServiceBEnabled,
		State:          memprovider.Enabled,
		DefaultVariant: "on",
		Variants: map[string]any{
			"on":  true,
			"off": false,
		}}

	FlagC = memprovider.InMemoryFlag{
		Key:            FlagServiceCEnabled,
		State:          memprovider.Enabled,
		DefaultVariant: "on",
		Variants: map[string]any{
			"on":  true,
			"off": false,
		}}

	FlagD = memprovider.InMemoryFlag{
		Key:            FlagServiceDEnabled,
		State:          memprovider.Enabled,
		DefaultVariant: "on",
		Variants: map[string]any{
			"on":  true,
			"off": false,
		}}
)

type Latch struct {
	count int
	mutex sync.Mutex
	cond  *sync.Cond
}

func NewLatch() *Latch {
	latch := &Latch{count: 1}
	latch.cond = sync.NewCond(&latch.mutex)
	return latch
}

func (l *Latch) Wait() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	for l.count > 0 {
		l.cond.Wait()
	}
}

func (l *Latch) CountDown() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.count--
	l.cond.Broadcast()
}

type Service struct {
	Flag string
}

func NewService(flag string) *Service {
	return &Service{
		Flag: flag,
	}
}

func (s *Service) Serve(ctx context.Context) error {
	features := openfeature.NewDefaultClient()

	if !features.Boolean(ctx, s.Flag, false, openfeature.TransactionContext(ctx)) {
		return fmt.Errorf("%s is false", s.Flag)
	}

	return nil
}
