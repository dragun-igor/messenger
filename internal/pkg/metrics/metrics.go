package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "messenger_service"

const (
	fieldMethodName = "grpc_method"
)

var requestTimeHist = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Subsystem: namespace,
	Name:      "request_duration_seconds",
	Help:      "Request duration per grpc method.",
}, []string{fieldMethodName})

var requestErrorsCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Subsystem: namespace,
	Name:      "request_errors_number",
	Help:      "Errors number method returns",
}, []string{fieldMethodName})
