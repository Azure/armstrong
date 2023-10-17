package provider

import (
	"github.com/ms-henglu/pal/rawlog"
	"github.com/ms-henglu/pal/types"
)

type Provider interface {
	IsTrafficTrace(l rawlog.RawLog) bool
	IsRequestTrace(l rawlog.RawLog) bool
	IsResponseTrace(l rawlog.RawLog) bool
	ParseTraffic(l rawlog.RawLog) (*types.RequestTrace, error)
	ParseRequest(l rawlog.RawLog) (*types.RequestTrace, error)
	ParseResponse(l rawlog.RawLog) (*types.RequestTrace, error)
}
