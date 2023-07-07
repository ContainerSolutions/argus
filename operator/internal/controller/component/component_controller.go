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

package component

import (
	"context"
	"time"

	res "github.com/ContainerSolutions/argus/operator/internal/component"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	"github.com/go-logr/logr"
)

// ComponentReconciler reconciles a Component object
type Reconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=argus.io,resources=components,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=argus.io,resources=components/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=argus.io,resources=components/finalizers,verbs=update

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("Component", req.NamespacedName)
	// Get Component
	Component := argusiov1alpha1.Component{}
	log.Info("Reconciling Component", "Component", Component.Name)
	err := r.Client.Get(ctx, req.NamespacedName, &Component)
	if apierrors.IsNotFound(err) {
		return ctrl.Result{}, nil
	} else if err != nil {
		log.Error(err, "could not get Component")
		return ctrl.Result{}, nil
	}
	originalRes := Component.DeepCopy()
	// Get ComponentControls with labels matching this Component
	ComponentControlList := argusiov1alpha1.ComponentControlList{}
	err = r.Client.List(ctx, &ComponentControlList, client.MatchingLabels{"argus.io/Component": Component.Name})
	if err != nil {
		log.Error(err, "could not list ComponentControls to update compliance status for %v", Component.Name)
		return ctrl.Result{}, err
	}
	res.UpdateControls(ComponentControlList, &Component)
	err = r.Client.Status().Patch(ctx, &Component, client.MergeFrom(originalRes))
	if err != nil {
		// Should we error here?
		log.Error(err, "could not update Component Controls")
	}
	err = res.UpdateChild(ctx, r.Client, &Component)
	if err != nil {
		// Should we error here?
		log.Error(err, "could not update parent child definition")
	}
	// Check Parents
	// TODO Make this a configuration
	return ctrl.Result{RequeueAfter: 1 * time.Hour}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&argusiov1alpha1.Component{}).
		Complete(r)
}
