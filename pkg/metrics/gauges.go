package metrics

import "github.com/prometheus/client_golang/prometheus"

var NodeUsage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "node_usage",
	Help: "Monitoring node usage",
}, []string{"node", "namespace"})
