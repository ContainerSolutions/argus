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

package componentassessment

import (
	"context"
	"fmt"
	"time"

	lib "github.com/ContainerSolutions/argus/operator/internal/componentassessment"
	"github.com/ContainerSolutions/argus/operator/internal/metrics"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	"github.com/go-logr/logr"
)

// ComponentAssessmentReconciler reconciles a ComponentAssessment object
type ComponentAssessmentReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=argus.io,resources=componentassessments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=argus.io,resources=componentassessments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=argus.io,resources=componentassessments/finalizers,verbs=update

func (r *ComponentAssessmentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("ComponentAssessment", req.NamespacedName)
	// Get Component
	res := argusiov1alpha1.ComponentAssessment{}
	err := r.Client.Get(ctx, req.NamespacedName, &res)
	if apierrors.IsNotFound(err) {
		return ctrl.Result{}, nil
	} else if err != nil {
		log.Error(err, "could not get Component")
		return ctrl.Result{}, nil
	}
	//log.Info("Reconciling ComponentAssessment", "ComponentAssessment", res.Name)
	attestations, err := lib.ListComponentAttestations(ctx, r.Client, &res)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not list ComponentAttestations: %w", err)
	}
	children, valid := lib.GetValidComponentAttestations(ctx, attestations)
	original := res.DeepCopy()
	res.Status.ComponentAttestations = children
	res.Status.TotalAttestations = len(children)
	res.Status.PassedAttestations = valid
	res.Status.RunAt = metav1.Now()
	labels := map[string]string{
		"Component":  res.Labels["argus.io/Component"],
		"Assessment": res.Labels["argus.io/Assessment"],
		"Control":    res.Labels["argus.io/Control"],
	}
	metrics.GetGaugeVec(metrics.AttestationTotalKey).With(labels).Set(float64(res.Status.TotalAttestations))
	metrics.GetGaugeVec(metrics.AttestationValidKey).With(labels).Set(float64(res.Status.PassedAttestations))
	err = r.Client.Status().Patch(ctx, &res, client.MergeFrom(original))
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not update Control status: %w", err)
	}
	// // Update ComponentControls metadata (force reconciliation)
	// list := argusiov1alpha1.ComponentControlList{}
	// err = r.Client.List(ctx, &list, client.MatchingLabels{"argus.io/Control": res.Labels["argus.io/Control"], "argus.io/Component": res.Labels["argus.io/Component"]})
	// if err != nil {
	// 	return ctrl.Result{}, fmt.Errorf("could not get ComponentControl: %w", err)
	// }
	// for _, resReq := range list.Items {
	// 	originalResReq := resReq.DeepCopy()
	// 	resReq.ObjectMeta.Annotations = map[string]string{fmt.Sprintf("%v.Assessment.argus.io/lastRun", res.Name): time.Now().String()}
	// 	err = r.Client.Patch(ctx, &resReq, client.MergeFrom(originalResReq))
	// 	if err != nil {
	// 		return ctrl.Result{}, fmt.Errorf("could not trigger ComponentControl reconciliation: %w", err)
	// 	}

	// }

	return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ComponentAssessmentReconciler) SetupWithManager(mgr ctrl.Manager, opts controller.Options) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&argusiov1alpha1.ComponentAssessment{}).
		WithOptions(opts).
		Complete(r)
}
