/*
Copyright 2021.

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

package controllers

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	ratelimitv1alpha1 "github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/service"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// RateLimitServiceReconciler reconciles a RateLimitService object
type RateLimitServiceReconciler struct {
	Client client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=ratelimit.zufardhiyaulhaq.com,resources=ratelimitservices,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ratelimit.zufardhiyaulhaq.com,resources=ratelimitservices/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ratelimit.zufardhiyaulhaq.com,resources=ratelimitservices/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the RateLimitService object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *RateLimitServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Start RateLimitService Reconciler")

	rateLimitService := &ratelimitv1alpha1.RateLimitService{}
	err := r.Client.Get(ctx, req.NamespacedName, rateLimitService)
	if err != nil {
		return ctrl.Result{}, nil
	}

	globalRateLimitConfigList := &ratelimitv1alpha1.GlobalRateLimitConfigList{}
	err = r.Client.List(ctx, globalRateLimitConfigList, &client.ListOptions{})
	if err != nil {
		return ctrl.Result{}, err
	}

	globalRateLimitList := &ratelimitv1alpha1.GlobalRateLimitList{}
	err = r.Client.List(ctx, globalRateLimitList, &client.ListOptions{})
	if err != nil {
		return ctrl.Result{}, err
	}

	ownGlobalRateLimitConfigList := []ratelimitv1alpha1.GlobalRateLimitConfig{}
	ownGlobalRateLimitList := []ratelimitv1alpha1.GlobalRateLimit{}

	// fetch all GlobalRateLimitConfig that refer the service to RateLimitService
	for _, globalRateLimitConfig := range globalRateLimitConfigList.Items {
		if globalRateLimitConfig.Spec.Ratelimit.Spec.Service.Type == "service" {
			if globalRateLimitConfig.Spec.Ratelimit.Spec.Service.Name == req.Name && globalRateLimitConfig.Namespace == req.Namespace {
				ownGlobalRateLimitConfigList = append(ownGlobalRateLimitConfigList, globalRateLimitConfig)
			}
		}
	}

	if len(ownGlobalRateLimitConfigList) == 0 {
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	// check all GlobalRateLimitConfig have the same ratelimit domain
	for a, configA := range globalRateLimitConfigList.Items {
		for b, configB := range globalRateLimitConfigList.Items {
			if a == b {
				continue
			}

			if configA.Spec.Ratelimit.Spec.Domain != configB.Spec.Ratelimit.Spec.Domain {
				return ctrl.Result{}, fmt.Errorf("globalRateLimitConfig use different domain")
			}
		}
	}

	// fetch all GlobalRateLimit that refer to GlobalRateLimitConfig
	for _, globalRateLimitConfig := range ownGlobalRateLimitConfigList {
		for _, globalRateLimit := range globalRateLimitList.Items {
			if globalRateLimit.Namespace == globalRateLimitConfig.Namespace && globalRateLimit.Spec.Config == globalRateLimitConfig.Name {
				ownGlobalRateLimitList = append(ownGlobalRateLimitList, globalRateLimit)
			}
		}
	}

	if len(ownGlobalRateLimitList) == 0 {
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	log.Info("Building Descriptors")
	descriptors, err := service.NewRateLimitDescriptor(ownGlobalRateLimitList)
	if err != nil {
		return ctrl.Result{}, err
	}

	svc, err := service.NewServiceBuilder().
		SetRateLimitService(*rateLimitService).
		Build()
	if err != nil {
		return ctrl.Result{}, err
	}

	config, err := service.NewRateLimitConfig(ownGlobalRateLimitConfigList[0].Spec.Ratelimit.Spec.Domain, descriptors)
	if err != nil {
		return ctrl.Result{}, err
	}
	configString, err := config.String()
	if err != nil {
		return ctrl.Result{}, err
	}

	configMapConfig, err := service.NewConfigBuilder().
		SetRateLimitService(*rateLimitService).
		SetConfig(configString).
		Build()
	if err != nil {
		return ctrl.Result{}, err
	}

	configMapEnv, err := service.NewEnvBuilder().
		SetRateLimitService(*rateLimitService).
		Build()
	if err != nil {
		return ctrl.Result{}, err
	}

	deployment, err := service.NewDeploymentBuilder().
		SetRateLimitService(*rateLimitService).
		Build()
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("set reference service")
	err = ctrl.SetControllerReference(rateLimitService, svc, r.Scheme)
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("set reference configmap config")
	err = ctrl.SetControllerReference(rateLimitService, configMapConfig, r.Scheme)
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("set reference configmap env")
	err = ctrl.SetControllerReference(rateLimitService, configMapEnv, r.Scheme)
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("set reference deployment")
	err = ctrl.SetControllerReference(rateLimitService, deployment, r.Scheme)
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("get service")
	createdSvc := &corev1.Service{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: svc.Name, Namespace: svc.Namespace}, createdSvc)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("create service")
			err = r.Client.Create(ctx, svc)
			if err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{Requeue: true}, nil
		} else {
			return ctrl.Result{}, err
		}
	}

	log.Info("get configmap config")
	createdConfigMapConfig := &corev1.ConfigMap{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: configMapConfig.Name, Namespace: configMapConfig.Namespace}, createdConfigMapConfig)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("create configmap config")
			err = r.Client.Create(ctx, configMapConfig)
			if err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{Requeue: true}, nil
		} else {
			return ctrl.Result{}, err
		}
	}

	log.Info("get configmap env")
	createdConfigMapEnv := &corev1.ConfigMap{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: configMapEnv.Name, Namespace: configMapEnv.Namespace}, createdConfigMapEnv)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("create configmap env")
			err = r.Client.Create(ctx, configMapEnv)
			if err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{Requeue: true}, nil
		} else {
			return ctrl.Result{}, err
		}
	}

	log.Info("get deployment")
	createdDeployment := &appsv1.Deployment{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, createdDeployment)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("create deployment")
			err = r.Client.Create(ctx, deployment)
			if err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{Requeue: true}, nil
		} else {
			return ctrl.Result{}, err
		}
	}

	if !equality.Semantic.DeepEqual(createdConfigMapConfig.Data, configMapConfig.Data) {
		createdConfigMapConfig.Data = configMapConfig.Data

		log.Info("update configmap config")
		err := r.Client.Update(ctx, createdConfigMapConfig, &client.UpdateOptions{})
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	if !equality.Semantic.DeepEqual(createdConfigMapEnv.Data, configMapEnv.Data) {
		createdConfigMapEnv.Data = configMapEnv.Data

		log.Info("update configmap env")
		err := r.Client.Update(ctx, createdConfigMapEnv, &client.UpdateOptions{})
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{RequeueAfter: 60 * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RateLimitServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ratelimitv1alpha1.RateLimitService{}).
		Complete(r)
}
