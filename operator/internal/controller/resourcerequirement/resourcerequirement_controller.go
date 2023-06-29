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

package resourcerequirement

import (
	"context"
	"fmt"
	"time"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	lib "github.com/ContainerSolutions/argus/operator/internal/resourcerequirement"
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ResourceRequirementReconciler reconciles a ResourceRequirement object
type ResourceRequirementReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=argus.io,resources=resourcerequirements,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=argus.io,resources=resourcerequirements/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=argus.io,resources=resourcerequirements/finalizers,verbs=update

func (r *ResourceRequirementReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("ClusterExternalSecret", req.NamespacedName)
	// Get Resource
	res := argusiov1alpha1.ResourceRequirement{}
	err := r.Client.Get(ctx, req.NamespacedName, &res)
	if apierrors.IsNotFound(err) {
		return ctrl.Result{}, nil
	} else if err != nil {
		log.Error(err, "could not get resource")
		return ctrl.Result{}, nil
	}
	implementations, valid, err := lib.GetValidResourceImplementations(ctx, r.Client, res)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not get resource implementations for requirement '%v': %w", res.Name, err)
	}
	original := res.DeepCopy()
	res.Status.ValidImplementations = valid
	res.Status.TotalImplementations = len(implementations)
	res.Status.ApplicableResourceImplementations = implementations
	res.Status.Status = "Not Implemented"
	if res.Status.TotalImplementations == res.Status.ValidImplementations && res.Status.TotalImplementations > 0 {
		res.Status.Status = "Implemented"
	}
	err = r.Client.Status().Patch(ctx, &res, client.MergeFrom(original))
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not update resourcerequirement status: %w", err)
	}
	return ctrl.Result{RequeueAfter: 4 * time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ResourceRequirementReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&argusiov1alpha1.ResourceRequirement{}).
		Complete(r)
}
