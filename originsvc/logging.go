package main

import (
	"time"

	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   OriginService
}

func (mw loggingMiddleware) ProcessAlert(sid string) (output string, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "process_alert",
			"input", sid,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.next.ProcessAlert(sid)
	return
}

func (mw loggingMiddleware) List(s string) (output map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "list",
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.next.List(s)
	return
}
