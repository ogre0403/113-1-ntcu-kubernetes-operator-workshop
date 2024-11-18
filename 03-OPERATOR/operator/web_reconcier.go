package main

import (
	"context"

	webv1 "operator/pkg/apis/myweb/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type WebReconciler struct {
	client client.Client
	scheme *runtime.Scheme
}

func (r *WebReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {

	log := log.FromContext(ctx)

	sample := &webv1.MyWeb{}
	err := r.client.Get(ctx, req.NamespacedName, sample)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("MyWeb resource not found. Ignoring since object must be deleted")
			return reconcile.Result{}, nil
		}
		log.Error(err, "Failed to get sample")
		return reconcile.Result{}, err
	}

	foundDeployment := &appsv1.Deployment{}
	err = r.client.Get(ctx, types.NamespacedName{Name: sample.Name, Namespace: sample.Namespace}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		dep := r.newDeployment(sample)
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return reconcile.Result{}, err
	}

	foundSvc := &corev1.Service{}
	err = r.client.Get(ctx, types.NamespacedName{Name: sample.Name, Namespace: sample.Namespace}, foundSvc)

	if err != nil && errors.IsNotFound(err) {
		svc := r.newService(sample)
		log.Info("Creating a new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
		err = r.client.Create(ctx, svc)
		if err != nil {
			log.Error(err, "Failed to create new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Service")
		return reconcile.Result{}, err
	}

	specPort := sample.Spec.NodePortNumber
	currentPort := int(foundSvc.Spec.Ports[0].NodePort)

	if specPort != 0 && int32(currentPort) != int32(specPort) {
		foundSvc.Spec.Ports[0].NodePort = int32(specPort)

		err = r.client.Update(ctx, foundSvc)
		if err != nil {
			log.Error(err, "Failed to update Service", "Service.Namespace", foundSvc.Namespace, "Service.Name", foundSvc.Name)
			return reconcile.Result{}, err
		}

		if sample.Status.Port != specPort {
			sample.Status.Port = specPort
			err = r.client.Status().Update(ctx, sample)
			if err != nil {
				log.Error(err, "Failed to update MyWeb status")
				return reconcile.Result{}, err
			}
		}

		return reconcile.Result{Requeue: true}, nil
	}

	// Update the MyWeb status with to true
	if !sample.Status.Completed {
		sample.Status.Completed = true
		err = r.client.Status().Update(ctx, sample)
		if err != nil {
			log.Error(err, "Failed to update MyWeb status")
			return reconcile.Result{}, err
		}
	}

	if sample.Status.Port == 0 || sample.Status.Port != specPort {
		sample.Status.Port = currentPort
		err = r.client.Status().Update(ctx, sample)
		if err != nil {
			log.Error(err, "Failed to update MyWeb status")
			return reconcile.Result{}, err
		}
	}

	log.Info("Exiting Reconcile")

	return reconcile.Result{}, nil
}

func (r *WebReconciler) newDeployment(s *webv1.MyWeb) *appsv1.Deployment {
	// take reference of a integer
	replicas := int32(1)

	deployment := &appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{

			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": s.Name},
			},

			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   s.Name,
					Labels: map[string]string{"app": s.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: s.Spec.Image,
						Name:  s.Name,
						Ports: []corev1.ContainerPort{{ContainerPort: 80}},
					}},
				},
			},
		},
	}

	deployment.Name = s.Name
	deployment.Namespace = s.Namespace

	ctrl.SetControllerReference(s, deployment, r.scheme)
	return deployment
}

func (r *WebReconciler) newService(s *webv1.MyWeb) *corev1.Service {

	svc := &corev1.Service{
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeNodePort,
			Selector: map[string]string{"app": s.Name},
			Ports: []corev1.ServicePort{
				{
					Protocol: corev1.ProtocolTCP,
					Port:     80,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 80,
					},
					// NodePort: int32(s.Spec.NodePortNumber),
				},
			},
		},
	}

	if s.Spec.NodePortNumber != 0 {
		svc.Spec.Ports[0].NodePort = int32(s.Spec.NodePortNumber)
	}

	svc.Name = s.Name
	svc.Namespace = s.Namespace

	ctrl.SetControllerReference(s, svc, r.scheme)
	return svc
}
