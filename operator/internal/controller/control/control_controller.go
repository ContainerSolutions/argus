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

package control

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	reqlib "github.com/ContainerSolutions/argus/operator/internal/control"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	"github.com/go-logr/logr"
)

// ControlReconciler reconciles a Control object
type ControlReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=argus.io,resources=controls,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=argus.io,resources=controls/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=argus.io,resources=controls/finalizers,verbs=update

func (r *ControlReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("Control", req.NamespacedName)
	// Get Component
	Control := argusiov1alpha1.Control{}
	err := r.Client.Get(ctx, req.NamespacedName, &Control)
	log.Info("Reconciling Control", "Control", Control.Name)
	if apierrors.IsNotFound(err) {
		return ctrl.Result{}, nil
	} else if err != nil {
		log.Error(err, "could not get Component")
		return ctrl.Result{}, nil
	}
	ComponentList := argusiov1alpha1.ComponentList{}
	err = r.Client.List(ctx, &ComponentList)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not list Components CR: %w", err)
	}
	Components := ComponentList.Items
	currentResReqs, err := reqlib.GetComponentControlsFromControl(ctx, r.Client, &Control)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not get Componentsrequiements for Control '%v': %w", Control.Name, err)
	}
	err = reqlib.LifecycleComponentControls(ctx, r.Client, Control.Spec.ApplicableComponentClasses, Components, currentResReqs)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not remove uneeded ComponentControls: %w", err)
	}
	children, err := reqlib.CreateOrUpdateComponentControls(ctx, r.Client, r.Scheme, &Control, Components)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not create ComponentControl for Control '%v': %w", Control.Name, err)
	}
	// Update Control Status
	original := Control.DeepCopy()
	Control.Status.Children = children
	ControlSpecbytes, err := json.Marshal(Control.Spec.Definition)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not marshal Control spec: %w", err)
	}
	hexSHA := sha512.Sum512(ControlSpecbytes)
	Control.Status.ControlHash = hex.EncodeToString(hexSHA[:])
	err = r.Client.Status().Patch(ctx, &Control, client.MergeFrom(original))
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not update Control status: %w", err)
	}
	return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ControlReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&argusiov1alpha1.Control{}).
		Complete(r)
}
