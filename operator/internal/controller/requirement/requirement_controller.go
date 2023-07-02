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

package requirement

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	reqlib "github.com/ContainerSolutions/argus/operator/internal/requirement"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	"github.com/go-logr/logr"
)

// RequirementReconciler reconciles a Requirement object
type RequirementReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=argus.io,resources=requirements,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=argus.io,resources=requirements/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=argus.io,resources=requirements/finalizers,verbs=update

func (r *RequirementReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("ClusterExternalSecret", req.NamespacedName)
	// Get Resource
	requirement := argusiov1alpha1.Requirement{}
	err := r.Client.Get(ctx, req.NamespacedName, &requirement)
	if apierrors.IsNotFound(err) {
		return ctrl.Result{}, nil
	} else if err != nil {
		log.Error(err, "could not get resource")
		return ctrl.Result{}, nil
	}
	resourceList := argusiov1alpha1.ResourceList{}
	err = r.Client.List(ctx, &resourceList)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not list Resources CR: %w", err)
	}
	resources := resourceList.Items
	currentResReqs, err := reqlib.GetResourceRequirementsFromRequirement(ctx, r.Client, &requirement)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not get resourcesrequiements for requirement '%v': %w", requirement.Name, err)
	}
	err = reqlib.LifecycleResourceRequirements(ctx, r.Client, requirement.Spec.ApplicableResourceClasses, resources, currentResReqs)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not remove uneeded resourcerequirements: %w", err)
	}
	children, err := reqlib.CreateOrUpdateResourceRequirements(ctx, r.Client, r.Scheme, &requirement, resources)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not create ResourceRequirement for requirement '%v': %w", requirement.Name, err)
	}
	// Update Requirement Status
	original := requirement.DeepCopy()
	requirement.Status.Children = children
	requirementSpecbytes, err := json.Marshal(requirement.Spec.Definition)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not marshal requirement spec: %w", err)
	}
	hexSHA := sha512.Sum512(requirementSpecbytes)
	requirement.Status.RequirementHash = hex.EncodeToString(hexSHA[:])
	err = r.Client.Status().Patch(ctx, &requirement, client.MergeFrom(original))
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not update requirement status: %w", err)
	}
	return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RequirementReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&argusiov1alpha1.Requirement{}).
		Complete(r)
}
