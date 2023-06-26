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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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
	err := r.Get(ctx, req.NamespacedName, &resource)
	if apierrors.IsNotFound(err) {
		return ctrl.Result{}, nil
	} else if err != nil {
		log.Error(err, "could not get resource")
		return ctrl.Result{}, nil
	}

	// Check Parents
	for _, parentName := range resource.Spec.Parents {
		parentResource := argusiov1alpha1.Resource{}
		namespacedName := types.NamespacedName{
			Name:      parentName,
			Namespace: resource.Namespace,
		}
		// Get ResourceRequirements with labels matching this Resource
		// TODO
		// Update Compliant status based on ResourceRequirement Status
		// TODO
		// Update parent adding current as a Child
		err := r.Get(ctx, namespacedName, &parentResource)
		if err != nil {
			log.Error(err, "could not get parent resource")
			return ctrl.Result{}, err
		}
		original := parentResource.DeepCopy()
		if parentResource.Status.Children == nil {
			parentResource.Status.Children = make(map[string]argusiov1alpha1.ResourceChild)
		}
		parentResource.Status.Children[resource.Name] = argusiov1alpha1.ResourceChild{
			Compliant: false,
		}
		parentResource.Status.CompliantChildren = 1
		parentResource.Status.ImplementedRequirements = 1
		parentResource.Status.TotalRequirements = 1
		parentResource.Status.TotalChildren = 1
		err = r.Status().Patch(ctx, &parentResource, client.MergeFrom(original))
		if err != nil {
			log.Error(err, "could not update parent resource children")
		}
	}
	// TODO Make this a configuration
	return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&argusiov1alpha1.Resource{}).
		Complete(r)
}
