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

package componentattestation

import (
	"context"
	"fmt"
	"time"

	lib "github.com/ContainerSolutions/argus/operator/internal/componentattestation"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	"github.com/go-logr/logr"
)

// ComponentAttestationReconciler reconciles a ComponentAttestation object
type ComponentAttestationReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=argus.io,resources=componentattestations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=argus.io,resources=componentattestations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=argus.io,resources=componentattestations/finalizers,verbs=update
//+kubebuilder:rbac:groups=argus.io,resources=attestationproviders,verbs=get;list;watch

func (r *ComponentAttestationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var err error
	log := r.Log.WithValues("ComponentAttestation", req.NamespacedName)
	// Get Component
	res := argusiov1alpha1.ComponentAttestation{}
	err = r.Client.Get(ctx, req.NamespacedName, &res)
	if apierrors.IsNotFound(err) {
		return ctrl.Result{}, nil
	} else if err != nil {
		log.Error(err, "could not get Component")
		return ctrl.Result{}, nil
	}
	log.Info("Reconciling ComponentAttestation", "ComponentAttestation", res.Name)
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
		return ctrl.Result{}, fmt.Errorf("could not update Control status: %w", err)
	}
	// Update ComponentAssessment Metadata (to force reconciliation)
	resImp := argusiov1alpha1.ComponentAssessment{}
	resImpName := fmt.Sprintf("%v-%v", res.Labels["argus.io/Assessment"], res.Labels["argus.io/Component"])
	err = r.Client.Get(ctx, types.NamespacedName{Namespace: res.Namespace, Name: resImpName}, &resImp)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not get ComponentAssessment: %w", err)
	}
	originalResImp := resImp.DeepCopy()
	resImp.ObjectMeta.Annotations = map[string]string{fmt.Sprintf("%v.attestation.argus.io/lastRun", res.Name): time.Now().String()}
	err = r.Client.Patch(ctx, &resImp, client.MergeFrom(originalResImp))
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not trigger ComponentAssessment reconciliation: %w", err)
	}
	return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ComponentAttestationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&argusiov1alpha1.ComponentAttestation{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(r)
}
