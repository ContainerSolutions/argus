package control

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

func GetComponentControlsFromControl(ctx context.Context, cl client.Client, Control *argusiov1alpha1.Control) (map[string]argusiov1alpha1.ComponentControl, error) {
	ComponentControlList := argusiov1alpha1.ComponentControlList{}
	err := cl.List(ctx, &ComponentControlList, client.MatchingLabels{"argus.io/Control": fmt.Sprintf("%v_%v", Control.Spec.Definition.Code, Control.Spec.Definition.Version)})
	if err != nil {
		return nil, fmt.Errorf("could not list ComponentControls: %w", err)
	}
	resReqs := make(map[string]argusiov1alpha1.ComponentControl)
	for _, ComponentControl := range ComponentControlList.Items {
		resReqs[ComponentControl.Name] = ComponentControl
	}
	return resReqs, nil
}

func LifecycleComponentControls(ctx context.Context, cl client.Client, classes []string, Components []argusiov1alpha1.Component, resReq map[string]argusiov1alpha1.ComponentControl) error {
	ComponentNames := []string{}
	for _, Component := range Components {
		ComponentNames = append(ComponentNames, Component.Name)
	}
	for _, ComponentControl := range resReq {
		refComponent, ok := ComponentControl.ObjectMeta.Labels["argus.io/Component"]
		if !ok {
			return fmt.Errorf("object '%v' does not contain expected label 'argus.io/Component'", ComponentControl.Name)
		}
		// If Component does not exist, it was deleted - we need to delete ComponentControl
		if !utils.Contains(ComponentNames, refComponent) {
			err := cl.Delete(ctx, &ComponentControl)
			if err != nil {
				return fmt.Errorf("could not delete ComponentControl '%v': %w", ComponentControl.Name, err)
			}
		}
		class, ok := ComponentControl.ObjectMeta.Labels["argus.io/Component-class"]
		if !ok {
			return fmt.Errorf("object '%v' does not contain expected label 'argus.io/Component-class'", ComponentControl.Name)
		}
		// If Control Class has changed, we need to delete ComponentControl
		if !utils.Contains(classes, class) {
			err := cl.Delete(ctx, &ComponentControl)
			if err != nil {
				return fmt.Errorf("could not delete ComponentControl '%v': %w", ComponentControl.Name, err)
			}
		}
		// If Component Class has changed, we need to delete ComponentControl
		for _, Component := range Components {
			if refComponent == Component.Name && !utils.Contains(Component.Spec.Classes, class) {
				err := cl.Delete(ctx, &ComponentControl)
				if err != nil {
					return fmt.Errorf("could not delete ComponentControl '%v': %w", ComponentControl.Name, err)
				}
			}
		}
	}
	return nil
}

func CreateOrUpdateComponentControls(ctx context.Context, cl client.Client, scheme *runtime.Scheme, req *argusiov1alpha1.Control, Components []argusiov1alpha1.Component) ([]argusiov1alpha1.NamespacedName, error) {
	all := []argusiov1alpha1.NamespacedName{}
	for _, class := range req.Spec.ApplicableComponentClasses {
		for _, Component := range Components {
			if utils.Contains(Component.Spec.Classes, class) {
				resReq := &argusiov1alpha1.ComponentControl{
					ObjectMeta: metav1.ObjectMeta{
						Name:      fmt.Sprintf("%v-%v", req.Name, Component.Name),
						Namespace: req.Namespace,
					},
				}
				emptyMutation := func() error {
					resReq.Spec.Definition = req.Spec.Definition
					resReq.Spec.RequiredAssessmentClasses = req.Spec.RequiredAssessmentClasses
					resReq.ObjectMeta.Labels = map[string]string{
						"argus.io/Control":         fmt.Sprintf("%v_%v", req.Spec.Definition.Code, req.Spec.Definition.Version),
						"argus.io/Component":       Component.Name,
						"argus.io/Component-class": class,
					}
					return nil
				}
				err := controllerutil.SetControllerReference(req, &resReq.ObjectMeta, scheme)
				if err != nil {
					return nil, fmt.Errorf("could not set controller reference for ComponentAssessment '%v': %w", resReq.Name, err)
				}
				_, err = ctrl.CreateOrUpdate(ctx, cl, resReq, emptyMutation)
				if err != nil {
					return nil, fmt.Errorf("could not create ComponentControl '%v': %w", resReq.Name, err)
				}
				child := argusiov1alpha1.NamespacedName{
					Name:      resReq.Name,
					Namespace: resReq.Namespace,
				}
				all = append(all, child)
			}
		}
	}
	return all, nil
}
