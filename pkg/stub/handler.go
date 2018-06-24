package stub

import (
	"context"

	"github.com/scholzj/dispatch-operator/pkg/apis/dispatch/v1alpha1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.Router:
		resourceExists := true
		d := &appsv1.Deployment{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Deployment",
				APIVersion: "apps/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      o.GetName(),
				Namespace: o.GetNamespace(),
			},
		}

		err := sdk.Get(d)
		if err != nil && errors.IsNotFound(err) {
			logrus.Infof("Deployment doesn't yet exist: %v", err)
			resourceExists = false
		} else if err != nil {
			logrus.Errorf("Failed to get router deployment: %v", err)
			return err
		}

		if resourceExists {
			logrus.Infof("Updating deployment %s", o.Name)
			err := sdk.Update(newRouter(o))
			if err != nil && !errors.IsAlreadyExists(err) {
				logrus.Errorf("Failed to update router deployment: %v", err)
				return err
			}
		} else {
			logrus.Infof("Creating deployment %s", o.Name)
			err := sdk.Create(newRouter(o))
			if err != nil && !errors.IsAlreadyExists(err) {
				logrus.Errorf("Failed to create router deployment: %v", err)
				return err
			}
		}
	}

	return nil
}

func newRouter(cr *v1alpha1.Router) *v1.Deployment {
	labels := map[string]string{
		"app": "qpid-dispatch",
		"name": cr.Name,
	}

	maxUnavailable := intstr.FromInt(1)
	maxSurge := intstr.FromInt(1)

	return &appsv1.Deployment {
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "Router",
				}),
			},
			Labels: labels,
		},
		Spec: v1.DeploymentSpec{
			Replicas:                &cr.Spec.Nodes,
			Selector:                &metav1.LabelSelector{MatchLabels: labels},
			Strategy:                v1.DeploymentStrategy{
				Type:          	"RollingUpdate",
				RollingUpdate: &v1.RollingUpdateDeployment{
					MaxUnavailable: &maxUnavailable,
					MaxSurge:       &maxSurge,
				},
			},
			Template:                corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:		labels,
				},
				Spec:       corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "router",
							Image:   "scholzj/qpid-dispatch:latest",
							Command: []string{"qdrouterd"},
							Ports: []corev1.ContainerPort{{
								Name:          "amqp",
								ContainerPort: int32(5672),
							}},
						},
					},
				},
			},
		},
	}
}