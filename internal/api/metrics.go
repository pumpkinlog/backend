package api

import "expvar"

var metric_endpoint_invocations *expvar.Map

func init() {
	metric_endpoint_invocations = expvar.NewMap("endpoint_invocations")
}
