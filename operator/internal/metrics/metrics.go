package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var AssessmentLabels = []string{"Component", "Assessment", "Control"}
var ControlLabels = []string{"Component", "Control"}
var ComponentLabels = []string{"Component"}

const (
	AttestationTotalKey = "attestations_total"
	AttestationValidKey = "attestations_valid"
	AssessmentTotalKey  = "Assessments_total"
	AssessmentValidKey  = "Assessments_valid"
	ControlTotalKey     = "Controls_total"
	ControlValidKey     = "Controls_valid"
)

var gaugeVecMetrics = map[string]*prometheus.GaugeVec{}

func SetUpMetrics() {
	// Only register once
	if len(gaugeVecMetrics) == 6 {
		return
	}
	// Obtain the prometheus metrics and register
	attestationsTotal := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "argus",
		Name:      "attestations_total",
		Help:      "Total number of Attestations",
	}, AssessmentLabels)
	attestationsValid := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "argus",
		Name:      "attestations_valid",
		Help:      "Total number of Attestations",
	}, AssessmentLabels)
	AssessmentsTotal := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "argus",
		Name:      "Assessments_total",
		Help:      "Total number of Assessments",
	}, ControlLabels)
	// Obtain the prometheus metrics and register
	AssessmentsValid := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "argus",
		Name:      "Assessments_valid",
		Help:      "Number of valid Assessments",
	}, ControlLabels)
	ControlsTotal := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "argus",
		Name:      "Controls_total",
		Help:      "Total number of Controls",
	}, ComponentLabels)
	ControlsValid := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "argus",
		Name:      "Controls_valid",
		Help:      "Number of valid Controls",
	}, ComponentLabels)
	metrics.Registry.MustRegister(
		attestationsTotal, attestationsValid,
		AssessmentsTotal, AssessmentsValid,
		ControlsTotal, ControlsValid)

	gaugeVecMetrics = map[string]*prometheus.GaugeVec{
		AttestationTotalKey: attestationsTotal,
		AttestationValidKey: attestationsValid,
		AssessmentTotalKey:  AssessmentsTotal,
		AssessmentValidKey:  AssessmentsValid,
		ControlTotalKey:     ControlsTotal,
		ControlValidKey:     ControlsValid,
	}
}

func GetGaugeVec(key string) *prometheus.GaugeVec {
	return gaugeVecMetrics[key]
}
