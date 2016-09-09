package main

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           OriginService
}

func (mw instrumentingMiddleware) ProcessAlert(sid string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "process_alert", "error", fmt.Sprint(err == nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.next.ProcessAlert(sid)
	return
}
