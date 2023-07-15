/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	"go.uber.org/zap/zapcore"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	"github.com/ContainerSolutions/argus/operator/internal/controller/assessment"
	"github.com/ContainerSolutions/argus/operator/internal/controller/attestation"
	"github.com/ContainerSolutions/argus/operator/internal/controller/component"
	"github.com/ContainerSolutions/argus/operator/internal/controller/componentassessment"
	"github.com/ContainerSolutions/argus/operator/internal/controller/componentattestation"
	"github.com/ContainerSolutions/argus/operator/internal/controller/componentcontrol"
	"github.com/ContainerSolutions/argus/operator/internal/controller/control"
	"github.com/ContainerSolutions/argus/operator/internal/metrics"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(argusiov1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var lvl zapcore.Level
	var enc zapcore.TimeEncoder
	metrics.SetUpMetrics()
	lvlErr := lvl.UnmarshalText([]byte("info"))
	if lvlErr != nil {
		setupLog.Error(lvlErr, "error unmarshalling loglevel")
		os.Exit(1)
	}
	encErr := enc.UnmarshalText([]byte("epoch"))
	if encErr != nil {
		setupLog.Error(encErr, "error unmarshalling timeEncoding")
		os.Exit(1)
	}

	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
		Level:       lvl,
		TimeEncoder: enc,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	config := ctrl.GetConfigOrDie()
	config.QPS = 500
	config.Burst = 1500

	mgr, err := ctrl.NewManager(config, ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "a09411a1.argus.io",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&attestation.AttestationReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    ctrl.Log,
	}).SetupWithManager(mgr, controller.Options{
		MaxConcurrentReconciles: 5,
	}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Attestation")
		os.Exit(1)
	}
	if err = (&assessment.AssessmentReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    ctrl.Log,
	}).SetupWithManager(mgr, controller.Options{
		MaxConcurrentReconciles: 5,
	}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Assessment")
		os.Exit(1)
	}
	if err = (&control.ControlReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    ctrl.Log,
	}).SetupWithManager(mgr, controller.Options{
		MaxConcurrentReconciles: 100,
	}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Control")
		os.Exit(1)
	}
	if err = (&component.Reconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    ctrl.Log,
	}).SetupWithManager(mgr, controller.Options{
		MaxConcurrentReconciles: 100,
	}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Component")
		os.Exit(1)
	}
	if err = (&componentcontrol.ComponentControlReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    ctrl.Log,
	}).SetupWithManager(mgr, controller.Options{
		MaxConcurrentReconciles: 100,
	}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ComponentControl")
		os.Exit(1)
	}
	if err = (&componentattestation.ComponentAttestationReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    ctrl.Log,
	}).SetupWithManager(mgr, controller.Options{
		MaxConcurrentReconciles: 100,
	}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ComponentAttestation")
		os.Exit(1)
	}
	if err = (&componentassessment.ComponentAssessmentReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    ctrl.Log,
	}).SetupWithManager(mgr, controller.Options{
		MaxConcurrentReconciles: 100,
	}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ComponentAssessment")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
