package main

import (
	"context"
	"time"
)

func NewWorker(ctx context.Context, period time.Duration, f func()) {
	tick := time.NewTicker(period)
	f()
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			f()
		}
	}
}
