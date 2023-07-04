package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var implementationLabels = []string{"resource", "implementation", "requirement"}
var requirementLabels = []string{"resource", "requirement"}
var resourceLabels = []string{"resource"}

const (
	AttestationTotalKey    = "attestations_total"
	AttestationValidKey    = "attestations_valid"
	ImplementationTotalKey = "implementations_total"
	ImplementationValidKey = "implementations_valid"
	RequirementTotalKey    = "requirements_total"
	RequirementValidKey    = "requirements_valid"
)

var gaugeVecMetrics = map[string]*prometheus.GaugeVec{}

func SetUpMetrics() {
	// Obtain the prometheus metrics and register
	attestationsTotal := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "argus",
		Name:      "attestations_total",
		Help:      "Total number of Attestations",
	}, implementationLabels)
	attestationsValid := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "argus",
		Name:      "attestations_valid",
		Help:      "Total number of Attestations",
	}, implementationLabels)
	implementationsTotal := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "argus",
		Name:      "implementations_total",
		Help:      "Total number of Implementations",
	}, requirementLabels)
	// Obtain the prometheus metrics and register
	implementationsValid := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "argus",
		Name:      "implementations_valid",
		Help:      "Number of valid Implementations",
	}, requirementLabels)
	requirementsTotal := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "argus",
		Name:      "requirements_total",
		Help:      "Total number of Requirements",
	}, resourceLabels)
	requirementsValid := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "argus",
		Name:      "requirements_valid",
		Help:      "Number of valid Requirements",
	}, resourceLabels)
	metrics.Registry.MustRegister(
		attestationsTotal, attestationsValid,
		implementationsTotal, implementationsValid,
		requirementsTotal, requirementsValid)

	gaugeVecMetrics = map[string]*prometheus.GaugeVec{
		AttestationTotalKey:    attestationsTotal,
		AttestationValidKey:    attestationsValid,
		ImplementationTotalKey: implementationsTotal,
		ImplementationValidKey: implementationsValid,
		RequirementTotalKey:    requirementsTotal,
		RequirementValidKey:    requirementsValid,
	}
}

func GetGaugeVec(key string) *prometheus.GaugeVec {
	return gaugeVecMetrics[key]
}
