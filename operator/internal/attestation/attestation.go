package attestation

import (
	"context"
	"fmt"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	"github.com/ContainerSolutions/argus/operator/internal/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func GetComponentAttestations(ctx context.Context, cl client.Client, res *argusiov1alpha1.Attestation) (map[string]argusiov1alpha1.ComponentAttestation, error) {
	ComponentAttestationList := argusiov1alpha1.ComponentAttestationList{}
	err := cl.List(ctx, &ComponentAttestationList, client.MatchingLabels{"argus.io/attestation": res.Name})
	if err != nil {
		return nil, fmt.Errorf("could not list ComponentAttestation: %w", err)
	}
	resReqs := make(map[string]argusiov1alpha1.ComponentAttestation)
	for _, item := range ComponentAttestationList.Items {
		resReqs[item.Name] = item
	}
	return resReqs, nil
}

func LifecycleComponentAttestations(ctx context.Context, cl client.Client, AssessmentRef string, Components []argusiov1alpha1.ComponentAssessment, items map[string]argusiov1alpha1.ComponentAttestation) error {
	ComponentNames := []string{}
	for _, Component := range Components {
		label, ok := Component.ObjectMeta.Labels["argus.io/Component"]
		if !ok {
			return fmt.Errorf("Component Assessment '%v' does not contain expected label 'argus.io/Component'", Component.Name)
		}
		ComponentNames = append(ComponentNames, label)
	}
	for _, item := range items {
		// if item does not belong anymore to the same Assessment ref, delete it (as the attestation was updated)
		if item.Labels["argus.io/Assessment"] != AssessmentRef {
			err := cl.Delete(ctx, &item)
			if err != nil {
				return fmt.Errorf("could not delete ComponentAttestation '%v': %w", item.Name, err)
			}
			continue
		}
		refComponent, ok := item.ObjectMeta.Labels["argus.io/Component"]
		if !ok {
			return fmt.Errorf("object '%v' does not contain expected label 'argus.io/Component'", item.Name)
		}
		// If Component does not exist, it was deleted - we need to delete ComponentControl
		if !utils.Contains(ComponentNames, refComponent) {
			err := cl.Delete(ctx, &item)
			if err != nil {
				return fmt.Errorf("could not delete ComponentAttestation '%v': %w", item.Name, err)
			}
		}
	}
	return nil

}

func CreateOrUpdateComponentAttestations(ctx context.Context, cl client.Client, scheme *runtime.Scheme, res *argusiov1alpha1.Attestation, Components []argusiov1alpha1.ComponentAssessment) ([]argusiov1alpha1.NamespacedName, error) {
	all := []argusiov1alpha1.NamespacedName{}
	for _, Component := range Components {
		AssessmentName, ok := Component.ObjectMeta.Labels["argus.io/Assessment"]
		if !ok {
			return nil, fmt.Errorf("Component Assessment '%v' does not contain expected label 'argus.io/Assessment'", Component.Name)
		}
		if res.Spec.AssessmentRef == AssessmentName {
			ComponentName, ok := Component.ObjectMeta.Labels["argus.io/Component"]
			if !ok {
				return nil, fmt.Errorf("Component Assessment '%v' does not contain expected label 'argus.io/Component'", Component.Name)
			}
			ControlName, ok := Component.ObjectMeta.Labels["argus.io/Control"]
			if !ok {
				return nil, fmt.Errorf("Component Assessment '%v' does not contain expected label 'argus.io/Control'", Component.Name)
			}
			resAtt := &argusiov1alpha1.ComponentAttestation{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%v-%v", res.Name, ComponentName),
					Namespace: res.Namespace,
				},
			}
			err := controllerutil.SetControllerReference(res, &resAtt.ObjectMeta, scheme)
			if err != nil {
				return nil, fmt.Errorf("could not set controller reference for ComponentAssessment '%v': %w", resAtt.Name, err)
			}
			emptyMutation := func() error {
				resAtt.Spec.ProviderRef = res.Spec.ProviderRef
				resAtt.ObjectMeta.Labels = map[string]string{
					"argus.io/Assessment":  res.Spec.AssessmentRef,
					"argus.io/attestation": res.Name,
					"argus.io/Component":   ComponentName,
					"argus.io/Control":     ControlName,
				}
				return nil
			}
			_, err = ctrl.CreateOrUpdate(ctx, cl, resAtt, emptyMutation)
			if err != nil {
				return nil, fmt.Errorf("could not create ComponentAttestation '%v': %w", resAtt.Name, err)
			}
			child := argusiov1alpha1.NamespacedName{
				Name:      resAtt.Name,
				Namespace: resAtt.Namespace,
			}
			all = append(all, child)
		}
	}
	return all, nil
}
