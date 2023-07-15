package component

import (
	"context"
	"fmt"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	"github.com/ContainerSolutions/argus/operator/internal/metrics"
	"github.com/hashicorp/go-multierror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func UpdateControls(ComponentControlList argusiov1alpha1.ComponentControlList, Component *argusiov1alpha1.Component) *argusiov1alpha1.Component {
	validControls := 0
	reqs := make(map[string]*argusiov1alpha1.ComponentControlCompliance)
	for _, ComponentControl := range ComponentControlList.Items {
		status := argusiov1alpha1.ComponentControlCompliance{}
		status.Implemented = false
		if (ComponentControl.Status.ValidAssessments == ComponentControl.Status.TotalAssessments) && ComponentControl.Status.TotalAssessments > 0 {
			status.Implemented = true
			validControls = validControls + 1
		}
		name := fmt.Sprintf("%v:%v", ComponentControl.Spec.Definition.Code, ComponentControl.Spec.Definition.Version)
		reqs[name] = &status
	}
	Component.Status.Controls = reqs
	Component.Status.RunAt = metav1.Now()
	Component.Status.TotalChildren = len(Component.Status.Children)
	compliantChildren := 0
	for _, child := range Component.Status.Children {
		if child.Compliant {
			compliantChildren = compliantChildren + 1
		}
	}
	Component.Status.CompliantChildren = compliantChildren
	Component.Status.TotalControls = len(ComponentControlList.Items)
	Component.Status.ImplementedControls = validControls
	labels := map[string]string{
		"Component": Component.Name,
	}
	metrics.GetGaugeVec(metrics.ControlTotalKey).With(labels).Set(float64(Component.Status.TotalControls))
	metrics.GetGaugeVec(metrics.ControlValidKey).With(labels).Set(float64(Component.Status.ImplementedControls))
	return Component
}

func UpdateChild(ctx context.Context, cl client.Client, Component *argusiov1alpha1.Component) error {
	var allErrors *multierror.Error
	for _, parentName := range Component.Spec.Parents {
		parentComponent := argusiov1alpha1.Component{}
		namespacedName := types.NamespacedName{
			Name:      parentName,
			Namespace: Component.Namespace,
		}
		err := cl.Get(ctx, namespacedName, &parentComponent)
		if err != nil {
			allErrors = multierror.Append(allErrors, fmt.Errorf("parent Component %v not found: %w", parentName, err))
			continue
		}
		original := parentComponent.DeepCopy()
		if parentComponent.Status.Children == nil {
			parentComponent.Status.Children = make(map[string]argusiov1alpha1.ComponentChild)
		}
		parentComponent.Status.Children[Component.Name] = argusiov1alpha1.ComponentChild{
			Compliant: Component.Status.TotalControls == Component.Status.ImplementedControls,
		}
		err = cl.Status().Patch(ctx, &parentComponent, client.MergeFrom(original))
		if err != nil {
			allErrors = multierror.Append(allErrors, fmt.Errorf("failed updating status for parent Component %v: %w", parentName, err))
			continue
		}
	}
	if allErrors != nil {
		return allErrors.ErrorOrNil()
	}
	return nil
}
