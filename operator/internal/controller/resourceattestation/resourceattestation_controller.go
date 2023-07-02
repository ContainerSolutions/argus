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

package resourceattestation

import (
	"context"
	"fmt"
	"time"

	lib "github.com/ContainerSolutions/argus/operator/internal/resourceattestation"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	"github.com/go-logr/logr"
)

// ResourceAttestationReconciler reconciles a ResourceAttestation object
type ResourceAttestationReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=argus.io,resources=resourceattestations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=argus.io,resources=resourceattestations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=argus.io,resources=resourceattestations/finalizers,verbs=update
//+kubebuilder:rbac:groups=argus.io,resources=attestationproviders,verbs=get;list;watch

func (r *ResourceAttestationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var err error
	log := r.Log.WithValues("ResourceAttestation", req.NamespacedName)
	// Get Resource
	res := argusiov1alpha1.ResourceAttestation{}
	err = r.Client.Get(ctx, req.NamespacedName, &res)
	if apierrors.IsNotFound(err) {
		return ctrl.Result{}, nil
	} else if err != nil {
		log.Error(err, "could not get resource")
		return ctrl.Result{}, nil
	}
	log.Info("Reconciling ResourceAttestation", "ResourceAttestation", res.Name)
	// Get Attestation Client
	attestationClient, err := lib.GetAttestationClient(ctx, r.Client, &res)
	if err != nil {
		return ctrl.Result{}, err
	}
	defer func() {
		e := attestationClient.Close()
		if e != nil && err == nil {
			err = fmt.Errorf("error closing client: %w", e)
		}
	}() // Prepare Call according to attestation provider logic
	result, err := attestationClient.Attest()
	if err != nil {
		return ctrl.Result{}, err
	}
	// Update Status
	original := res.DeepCopy()
	res.Status.Result = result
	res.Status.Status = "True"
	err = r.Client.Status().Patch(ctx, &res, client.MergeFrom(original))
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not update requirement status: %w", err)
	}
	return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ResourceAttestationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&argusiov1alpha1.ResourceAttestation{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(r)
}
