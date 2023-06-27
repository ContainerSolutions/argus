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

package resource

import (
	"context"
	"time"

	res "github.com/ContainerSolutions/argus/operator/internal/resource"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	"github.com/go-logr/logr"
)

// ResourceReconciler reconciles a Resource object
type Reconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=argus.io,resources=resources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=argus.io,resources=resources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=argus.io,resources=resources/finalizers,verbs=update

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("ClusterExternalSecret", req.NamespacedName)
	// Get Resource
	resource := argusiov1alpha1.Resource{}
	err := r.Client.Get(ctx, req.NamespacedName, &resource)
	if apierrors.IsNotFound(err) {
		return ctrl.Result{}, nil
	} else if err != nil {
		log.Error(err, "could not get resource")
		return ctrl.Result{}, nil
	}
	originalRes := resource.DeepCopy()
	// Get ResourceRequirements with labels matching this Resource
	resourceRequirementList := argusiov1alpha1.ResourceRequirementList{}
	err = r.Client.List(ctx, &resourceRequirementList, client.MatchingLabels{"argus.io/resource": resource.Name})
	if err != nil {
		log.Error(err, "could not list ResourceRequirements to update compliance status for %v", resource.Name)
		return ctrl.Result{}, err
	}
	res.UpdateRequirements(resourceRequirementList, &resource)
	err = r.Client.Status().Patch(ctx, &resource, client.MergeFrom(originalRes))
	if err != nil {
		// Should we error here?
		log.Error(err, "could not update resource Requirements")
	}
	err = res.UpdateChild(ctx, r.Client, &resource)
	if err != nil {
		// Should we error here?
		log.Error(err, "could not update parent child definition")
	}
	// Check Parents
	// TODO Make this a configuration
	return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&argusiov1alpha1.Resource{}).
		Complete(r)
}
