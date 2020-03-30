package request

import (
	"github.com/amazing-gao/gorequest"
	"github.com/boxgo/logger"
)

func logBefore(agent *gorequest.SuperAgent) {
	if GlobalOptions.UserAgent != "" {
		agent.AppendHeader("user-agent", GlobalOptions.UserAgent)
	}

	if GlobalOptions.ShowLog {
		curl, _ := agent.AsCurlCommand()
		logger.Trace(agent.Context).Infow("request_start", "curl", curl)
	}
}

func logAfter(agent *gorequest.SuperAgent, resp gorequest.Response, body []byte, errs []error) {
	if GlobalOptions.ShowLog {
		curl, _ := agent.AsCurlCommand()
		logger.Trace(agent.Context).Infow("request_end", "curl", curl, "errs", errs, "resp.status", resp.StatusCode, "body", string(body[:]), "resp.header", resp.Header)
	}
}
