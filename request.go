package request

import (
	"context"
	"time"

	"github.com/BiteBit/gorequest"
	"github.com/boxgo/box/minibox"
	"github.com/boxgo/metrics"
)

type (
	// Options http request tool
	Options struct {
		Timeout   int64  `config:"timeout" help:"Timeout millsecond, default 10s"`
		UserAgent string `config:"userAgent" help:"Client User-Agent"`
		ShowLog   bool   `config:"showLog" help:"Show request log"`
		Metrics   bool   `config:"metrics" help:"default is false"`

		app     minibox.App
		metrics metrics.Metrics
	}
)

var (
	// GlobalOptions global request options
	GlobalOptions = &Options{}

	befores []gorequest.Before
	afters  []gorequest.After
)

// Name config prefix name
func (opts *Options) Name() string {
	return "request"
}

// Exts 获取app信息
func (opts *Options) Exts() []minibox.MiniBox {
	return []minibox.MiniBox{&opts.app, &opts.metrics}
}

func (opts *Options) ConfigWillLoad(context.Context) {

}

func (opts *Options) ConfigDidLoad(context.Context) {
	if opts.UserAgent == "" {
		opts.UserAgent = opts.app.AppName
	}
}

// NewTraceRequest new a trace request
func NewTraceRequest(ctx context.Context) *gorequest.SuperAgent {
	agent := gorequest.NewWithContext(ctx)

	setup(agent)

	return agent
}

// UseBefore global use before
func UseBefore(bs ...gorequest.Before) {
	befores = append(befores, bs...)
}

// UseAfter global use after
func UseAfter(as ...gorequest.After) {
	afters = append(afters, as...)
}

func setup(agent *gorequest.SuperAgent) {
	timeout := time.Second * 10
	if GlobalOptions.Timeout != 0 {
		timeout = time.Duration(GlobalOptions.Timeout * int64(time.Millisecond))
	}

	agent.Timeout(timeout)

	agent.UseBefore(logBefore)
	agent.UseBefore(metricsBefore)
	agent.UseBefore(befores...)

	agent.UseAfter(logAfter)
	agent.UseAfter(metricsAfter)
	agent.UseAfter(afters...)
}
