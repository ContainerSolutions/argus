package assessment

import (
	"context"
	"fmt"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func GetComponentAssessments(ctx context.Context, cl client.Client, res *argusiov1alpha1.Assessment) (map[string]argusiov1alpha1.ComponentAssessment, error) {
	ComponentAssessmentList := argusiov1alpha1.ComponentAssessmentList{}
	err := cl.List(ctx, &ComponentAssessmentList, client.MatchingLabels{"argus.io/Assessment": res.Name})
	if err != nil {
		return nil, fmt.Errorf("could not list ComponentAssessment: %w", err)
	}
	resReqs := make(map[string]argusiov1alpha1.ComponentAssessment)
	for _, item := range ComponentAssessmentList.Items {
		resReqs[item.Name] = item
	}
	return resReqs, nil
}

func BuildComponentAssessmentList(ctx context.Context, res *argusiov1alpha1.Assessment, Components []argusiov1alpha1.Component) (map[string]argusiov1alpha1.ComponentAssessment, error) {
	items := map[string]argusiov1alpha1.ComponentAssessment{}
	// Treat Cascading policy. In order to do that, we need to add every child which this Assessment targets.
	ComponentNameList := []string{}
	ComponentMap := make(map[string]argusiov1alpha1.Component)
	for _, Component := range Components {
		ComponentMap[Component.Name] = Component
	}
	for _, Component := range Components {
		for _, refs := range res.Spec.ComponentRef {
			// If this Component is actually targetted by this Assessment
			if refs.Name == Component.Name && refs.Namespace == Component.Namespace {
				ComponentNameList = append(ComponentNameList, Component.Name)
				if res.Spec.CascadePolicy == argusiov1alpha1.CascadingPolicyCascade {
					for childName := range Component.Status.Children {
						// Add children
						ComponentNameList = append(ComponentNameList, childName)
					}
				}
			}

		}
	}
	for _, ComponentName := range ComponentNameList {
		Component := ComponentMap[ComponentName]
		resImp := argusiov1alpha1.ComponentAssessment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%v-%v", res.Name, Component.Name),
				Namespace: res.Namespace,
			},
		}
		items[resImp.Name] = resImp
	}
	return items, nil
}
func LifecycleComponentAssessments(ctx context.Context, cl client.Client, new, old map[string]argusiov1alpha1.ComponentAssessment) error {
	for name := range old {
		if _, ok := new[name]; !ok {
			del := old[name]
			err := cl.Delete(ctx, &del)
			if err != nil {
				return fmt.Errorf("could not delete ComponentAssessment '%v': %w", old[name].Name, err)
			}
		}
	}
	return nil
}

func CreateOrUpdateComponentAssessments(ctx context.Context, cl client.Client, scheme *runtime.Scheme, res *argusiov1alpha1.Assessment, Components []argusiov1alpha1.Component) ([]argusiov1alpha1.NamespacedName, error) {
	all := []argusiov1alpha1.NamespacedName{}
	// Treat Cascading policy. In order to do that, we need to add every child which this Assessment targets.
	ComponentNameList := []string{}
	ComponentMap := make(map[string]argusiov1alpha1.Component)
	for _, Component := range Components {
		ComponentMap[Component.Name] = Component
	}
	for _, Component := range Components {
		for _, refs := range res.Spec.ComponentRef {
			// If this Component is actually targetted by this Assessment
			if refs.Name == Component.Name && refs.Namespace == Component.Namespace {
				ComponentNameList = append(ComponentNameList, Component.Name)
				if res.Spec.CascadePolicy == argusiov1alpha1.CascadingPolicyCascade {
					for childName := range Component.Status.Children {
						// Add children
						ComponentNameList = append(ComponentNameList, childName)
					}
				}
			}

		}
	}
	for _, ComponentName := range ComponentNameList {
		Component := ComponentMap[ComponentName]
		resImp := &argusiov1alpha1.ComponentAssessment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%v-%v", res.Name, Component.Name),
				Namespace: res.Namespace,
			},
		}
		emptyMutation := func() error {
			resImp.Spec.ControlRef = res.Spec.ControlRef
			resImp.Spec.Class = res.Spec.Class
			resImp.ObjectMeta.Labels = map[string]string{
				"argus.io/Assessment": res.Name,
				"argus.io/Component":  Component.Name,
				"argus.io/Control":    fmt.Sprintf("%v_%v", res.Spec.ControlRef.Code, res.Spec.ControlRef.Version),
			}
			return nil
		}
		err := controllerutil.SetControllerReference(res, &resImp.ObjectMeta, scheme)
		if err != nil {
			return nil, fmt.Errorf("could not set controller reference for ComponentAssessment '%v': %w", resImp.Name, err)
		}
		_, err = ctrl.CreateOrUpdate(ctx, cl, resImp, emptyMutation)
		if err != nil {
			return nil, fmt.Errorf("could not create ComponentAssessment '%v': %w", resImp.Name, err)
		}
		child := argusiov1alpha1.NamespacedName{
			Name:      resImp.Name,
			Namespace: resImp.Namespace,
		}
		all = append(all, child)
	}
	return all, nil
}
