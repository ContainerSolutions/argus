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

package implementation

import (
	"context"
	"fmt"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	lib "github.com/ContainerSolutions/argus/operator/internal/implementation"
	"github.com/go-logr/logr"
)

// ImplementationReconciler reconciles a Implementation object
type ImplementationReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=argus.io,resources=implementations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=argus.io,resources=implementations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=argus.io,resources=implementations/finalizers,verbs=update

func (r *ImplementationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("Implementation", req.NamespacedName)
	// Get Resource
	res := argusiov1alpha1.Implementation{}
	err := r.Client.Get(ctx, req.NamespacedName, &res)
	if apierrors.IsNotFound(err) {
		return ctrl.Result{}, nil
	} else if err != nil {
		log.Error(err, "could not get resource")
		return ctrl.Result{}, nil
	}
	log.Info("Reconciling Implementation", "Implementation", res.Name)
	resourceList := argusiov1alpha1.ResourceList{}
	err = r.Client.List(ctx, &resourceList)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not list Resources CR: %w", err)
	}
	resources := resourceList.Items
	currentImplementations, err := lib.GetResourceImplementations(ctx, r.Client, &res)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not get ResourceImplementation for requirement '%v': %w", res.Name, err)
	}
	// TODO - If Requirement Resource Class does not match current Resources, no need to create new ResourceImplementations
	newList, err := lib.BuildResourceImplementationList(ctx, &res, resources)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not build list : %w", err)
	}
	err = lib.LifecycleResourceImplementations(ctx, r.Client, newList, currentImplementations)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not remove uneeded resourcerequirements: %w", err)
	}
	children, err := lib.CreateOrUpdateResourceImplementations(ctx, r.Client, r.Scheme, &res, resources)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not create ResourceImplementation for requirement '%v': %w", res.Name, err)
	}
	// Update Requirement Status
	original := res.DeepCopy()
	res.Status.Children = children
	err = r.Client.Status().Patch(ctx, &res, client.MergeFrom(original))
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not update requirement status: %w", err)
	}
	return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ImplementationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&argusiov1alpha1.Implementation{}).
		Complete(r)
}
