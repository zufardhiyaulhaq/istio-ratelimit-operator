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

	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/client/istio"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/global/config"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/utils"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	funk "github.com/thoas/go-funk"
	ratelimitv1alpha1 "github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// GlobalRateLimitConfigReconciler reconciles a GlobalRateLimitConfig object
type GlobalRateLimitConfigReconciler struct {
	client.Client
	IstioClient istio.ClientInterface
	Scheme      *runtime.Scheme
}

//+kubebuilder:rbac:groups=ratelimit.zufardhiyaulhaq.com,resources=globalratelimitconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ratelimit.zufardhiyaulhaq.com,resources=globalratelimitconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ratelimit.zufardhiyaulhaq.com,resources=globalratelimitconfigs/finalizers,verbs=update

func (r *GlobalRateLimitConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Start GlobalRateLimitConfig Reconciler")

	globalRateLimitConfig := &ratelimitv1alpha1.GlobalRateLimitConfig{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, globalRateLimitConfig)
	if err != nil {
		return ctrl.Result{}, nil
	}

	log.Info("Build envoyfilters")
	envoyFilters, err := config.NewConfigBuilder().
		SetConfig(*globalRateLimitConfig).
		Build()
	if err != nil {
		return ctrl.Result{}, err
	}

	if len(envoyFilters) == 0 {
		return ctrl.Result{}, fmt.Errorf("empty envoyfilter from builder")
	}

	allVersionEnvoyFilterNames := utils.BuildEnvoyFilterNamesAllVersion(globalRateLimitConfig.Name)
	EnvoyFilterNames := utils.BuildEnvoyFilterNames(globalRateLimitConfig.Name, globalRateLimitConfig.Spec.Selector.IstioVersion)

	deleteEnvoyFilters, _ := funk.DifferenceString(allVersionEnvoyFilterNames, EnvoyFilterNames)
	for _, deleteEnvoyFilter := range deleteEnvoyFilters {
		_, err := r.IstioClient.GetEnvoyFilter(ctx, globalRateLimitConfig.Namespace, deleteEnvoyFilter)
		if err != nil {
			continue
		}

		log.Info("delete unused envoyfilter")
		err = r.IstioClient.DeleteEnvoyFilter(ctx, globalRateLimitConfig.Namespace, deleteEnvoyFilter)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	log.Info("create envoyfilters")
	for _, envoyFilter := range envoyFilters {
		log.Info("set reference envoyfilter")
		if err := controllerutil.SetControllerReference(globalRateLimitConfig, envoyFilter, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}

		log.Info("get envoyfilter")
		createdEnvoyFilter, err := r.IstioClient.GetEnvoyFilter(ctx, envoyFilter.Namespace, envoyFilter.Name)
		if err != nil {
			if errors.IsNotFound(err) {
				log.Info("create envoyfilter")
				_, err := r.IstioClient.CreateEnvoyFilter(ctx, envoyFilter.Namespace, envoyFilter)
				if err != nil {
					return ctrl.Result{}, err
				}

				return ctrl.Result{Requeue: true}, nil
			} else {
				return ctrl.Result{}, err
			}
		}

		if !equality.Semantic.DeepEqual(createdEnvoyFilter.Spec, envoyFilter.Spec) {
			createdEnvoyFilter.Spec = envoyFilter.Spec

			log.Info("update envoyfilter")
			r.IstioClient.UpdateEnvoyFilter(ctx, envoyFilter.Namespace, createdEnvoyFilter)
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GlobalRateLimitConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ratelimitv1alpha1.GlobalRateLimitConfig{}).
		Complete(r)
}
