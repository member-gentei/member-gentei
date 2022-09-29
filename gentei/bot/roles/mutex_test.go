package roles

import (
	"context"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

func TestMutexErrGroup(t *testing.T) {
	var (
		mutex   = NewDefaultMapRWMutex()
		eg, ctx = errgroup.WithContext(context.Background())
	)
	const key = "test-key"
	eg.SetLimit(4)
	for i := 0; i < 4; i++ {
		eg.Go(func() error {
			m := mutex.GetOrCreate(key)
			m.Lock()
			defer m.Unlock()
			time.Sleep(time.Millisecond * 5)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		t.Fatalf("error returned from errGroup: %+v", err)
	}
	<-ctx.Done()
}
