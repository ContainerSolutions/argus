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

package componentcontrol

import (
	"context"
	"fmt"
	"time"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	lib "github.com/ContainerSolutions/argus/operator/internal/componentcontrol"
	"github.com/ContainerSolutions/argus/operator/internal/metrics"
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ComponentControlReconciler reconciles a ComponentControl object
type ComponentControlReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=argus.io,resources=ComponentControls,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=argus.io,resources=ComponentControls/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=argus.io,resources=ComponentControls/finalizers,verbs=update

func (r *ComponentControlReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("ComponentControl", req.NamespacedName)
	// Get Component
	res := argusiov1alpha1.ComponentControl{}
	err := r.Client.Get(ctx, req.NamespacedName, &res)
	if apierrors.IsNotFound(err) {
		return ctrl.Result{}, nil
	} else if err != nil {
		log.Error(err, "could not get Component")
		return ctrl.Result{}, nil
	}
	log.Info("Reconciling ComponentControl", "ComponentControl", res.Name)
	Assessments, valid, err := lib.GetValidComponentAssessments(ctx, r.Client, res)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not get Component Assessments for Control '%v': %w", res.Name, err)
	}
	original := res.DeepCopy()
	res.Status.ValidAssessments = valid
	res.Status.TotalAssessments = len(Assessments)
	res.Status.ApplicableComponentAssessments = Assessments
	res.Status.Status = "Not Implemented"
	res.Status.RunAt = metav1.Now()
	labels := map[string]string{
		"Component": res.Labels["argus.io/Component"],
		"Control":   res.Labels["argus.io/Control"],
	}
	metrics.GetGaugeVec(metrics.AssessmentTotalKey).With(labels).Set(float64(res.Status.TotalAssessments))
	metrics.GetGaugeVec(metrics.AssessmentValidKey).With(labels).Set(float64(res.Status.ValidAssessments))

	if res.Status.TotalAssessments == res.Status.ValidAssessments && res.Status.TotalAssessments > 0 {
		res.Status.Status = "Implemented"
	}
	err = r.Client.Status().Patch(ctx, &res, client.MergeFrom(original))
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not update ComponentControl status: %w", err)
	}
	// Update Component metadata (force reconciliation)
	list := argusiov1alpha1.Component{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: res.Labels["argus.io/Component"], Namespace: res.Namespace}, &list)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not get Component: %w", err)
	}
	originalResReq := list.DeepCopy()
	list.ObjectMeta.Annotations = map[string]string{fmt.Sprintf("%v.Control.argus.io/lastRun", res.Name): time.Now().String()}
	err = r.Client.Patch(ctx, &list, client.MergeFrom(originalResReq))
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not trigger Component reconciliation: %w", err)
	}
	return ctrl.Result{RequeueAfter: 1 * time.Hour}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ComponentControlReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&argusiov1alpha1.ComponentControl{}).
		Complete(r)
}
