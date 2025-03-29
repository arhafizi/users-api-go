package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var DbCall = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "db_calls_total",
	Help: "Number of calls to the database",
},
	[]string{"type_name", "operation_name", "status"})

var TotalReq = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "requests_total",
	Help: "Counting the total number of requests handled",
})
