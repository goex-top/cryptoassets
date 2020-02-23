package main

import (
	"context"
	"time"
)

func NewWorker(ctx context.Context, period time.Duration) {
	tick := time.NewTicker(period)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:

		}
	}
}
