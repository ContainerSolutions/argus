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

package resourceimplementation

import (
	"context"
	"fmt"
	"time"

	"github.com/ContainerSolutions/argus/operator/internal/metrics"
	lib "github.com/ContainerSolutions/argus/operator/internal/resourceimplementation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	"github.com/go-logr/logr"
)

// ResourceImplementationReconciler reconciles a ResourceImplementation object
type ResourceImplementationReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=argus.io,resources=resourceimplementations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=argus.io,resources=resourceimplementations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=argus.io,resources=resourceimplementations/finalizers,verbs=update

func (r *ResourceImplementationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("ResourceImplementation", req.NamespacedName)
	// Get Resource
	res := argusiov1alpha1.ResourceImplementation{}
	err := r.Client.Get(ctx, req.NamespacedName, &res)
	if apierrors.IsNotFound(err) {
		return ctrl.Result{}, nil
	} else if err != nil {
		log.Error(err, "could not get resource")
		return ctrl.Result{}, nil
	}
	log.Info("Reconciling ResourceImplementation", "ResourceImplementation", res.Name)
	attestations, err := lib.ListResourceAttestations(ctx, r.Client, &res)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not list ResourceAttestations: %w", err)
	}
	children, valid := lib.GetValidResourceAttestations(ctx, attestations)
	original := res.DeepCopy()
	res.Status.ResourceAttestations = children
	res.Status.TotalAttestations = len(children)
	res.Status.PassedAttestations = valid
	res.Status.RunAt = metav1.Now()
	labels := map[string]string{
		"resource":       res.Labels["argus.io/resource"],
		"implementation": res.Labels["argus.io/implementation"],
		"requirement":    res.Labels["argus.io/requirement"],
	}
	metrics.GetGaugeVec(metrics.AttestationTotalKey).With(labels).Set(float64(res.Status.TotalAttestations))
	metrics.GetGaugeVec(metrics.AttestationValidKey).With(labels).Set(float64(res.Status.PassedAttestations))
	err = r.Client.Status().Patch(ctx, &res, client.MergeFrom(original))
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not update requirement status: %w", err)
	}
	// Update ResourceRequirements metadata (force reconciliation)
	list := argusiov1alpha1.ResourceRequirementList{}
	err = r.Client.List(ctx, &list, client.MatchingLabels{"argus.io/requirement": res.Labels["argus.io/requirement"], "argus.io/resource": res.Labels["argus.io/resource"]})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not get resourcerequirement: %w", err)
	}
	for _, resReq := range list.Items {
		originalResReq := resReq.DeepCopy()
		resReq.ObjectMeta.Annotations = map[string]string{fmt.Sprintf("%v.implementation.argus.io/lastRun", res.Name): time.Now().String()}
		err = r.Client.Patch(ctx, &resReq, client.MergeFrom(originalResReq))
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("could not trigger resourcerequirement reconciliation: %w", err)
		}

	}

	return ctrl.Result{RequeueAfter: 1 * time.Hour}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ResourceImplementationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&argusiov1alpha1.ResourceImplementation{}).
		Complete(r)
}
