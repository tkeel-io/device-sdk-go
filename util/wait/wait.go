package wait

import (
    "context"
    "time"
)

func Every(f func(), period time.Duration, stopCh <-chan struct{}) *time.Ticker {
    ticker := time.NewTicker(period)

    go func() {
        for {
            select {
            case <-ticker.C:
                f()
            case <-stopCh:
                return
            }
        }
    }()

    return ticker
}

func EveryWithContext(ctx context.Context, f func(ctx context.Context), period time.Duration) *time.Ticker {
    return Every(func() { f(ctx) }, period, ctx.Done())
}
