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

	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/local/config"
	"github.com/zufardhiyaulhaq/istio-ratelimit-operator/pkg/utils"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	funk "github.com/thoas/go-funk"
	ratelimitv1alpha1 "github.com/zufardhiyaulhaq/istio-ratelimit-operator/api/v1alpha1"
	clientnetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
	ctrl "sigs.k8s.io/controller-runtime"
)

// LocalRateLimitConfigReconciler reconciles a LocalRateLimitConfig object
type LocalRateLimitConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=ratelimit.zufardhiyaulhaq.com,resources=localratelimitconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ratelimit.zufardhiyaulhaq.com,resources=localratelimitconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ratelimit.zufardhiyaulhaq.com,resources=localratelimitconfigs/finalizers,verbs=update

func (r *LocalRateLimitConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Start LocalRateLimitConfig Reconciler")

	localRateLimitConfig := &ratelimitv1alpha1.LocalRateLimitConfig{}
	err := r.Client.Get(ctx, req.NamespacedName, localRateLimitConfig)
	if err != nil {
		return ctrl.Result{}, nil
	}

	log.Info("Build localratelimitconfig envoyfilters")
	envoyFilters, err := config.NewConfigBuilder().
		SetConfig(*localRateLimitConfig).
		Build()
	if err != nil {
		return ctrl.Result{}, err
	}

	if len(envoyFilters) == 0 {
		return ctrl.Result{}, fmt.Errorf("empty localratelimitconfig envoyfilter from builder")
	}

	allVersionEnvoyFilterNames := utils.BuildEnvoyFilterNamesAllVersion(localRateLimitConfig.Name)
	EnvoyFilterNames := utils.BuildEnvoyFilterNames(localRateLimitConfig.Name, localRateLimitConfig.Spec.Selector.IstioVersion)

	deleteEnvoyFilters, _ := funk.DifferenceString(allVersionEnvoyFilterNames, EnvoyFilterNames)
	for _, deleteEnvoyFilterName := range deleteEnvoyFilters {
		deleteEnvoyFilter := &clientnetworking.EnvoyFilter{}
		err := r.Client.Get(ctx, types.NamespacedName{Name: deleteEnvoyFilterName, Namespace: req.Namespace}, deleteEnvoyFilter)
		if err != nil {
			continue
		}

		log.Info("delete unused localratelimitconfig envoyfilter")
		err = r.Client.Delete(ctx, deleteEnvoyFilter, &client.DeleteOptions{})
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	log.Info("create localratelimitconfig envoyfilters")
	for _, envoyFilter := range envoyFilters {
		err := ctrl.SetControllerReference(localRateLimitConfig, envoyFilter, r.Scheme)
		if err != nil {
			return ctrl.Result{}, err
		}

		log.Info("get localratelimitconfig envoyfilter")
		createdEnvoyFilter := &clientnetworking.EnvoyFilter{}
		err = r.Client.Get(ctx, types.NamespacedName{Name: envoyFilter.Name, Namespace: envoyFilter.Namespace}, createdEnvoyFilter)
		if err != nil {
			if errors.IsNotFound(err) {
				log.Info("create localratelimitconfig envoyfilter")
				err := r.Client.Create(ctx, envoyFilter, &client.CreateOptions{})
				if err != nil {
					return ctrl.Result{}, err
				}

				return ctrl.Result{RequeueAfter: 60 * time.Second}, nil
			} else {
				return ctrl.Result{}, err
			}
		}

		if !equality.Semantic.DeepEqual(createdEnvoyFilter.Spec, envoyFilter.Spec) {
			createdEnvoyFilter.Spec = envoyFilter.Spec

			log.Info("update localratelimitconfig envoyfilter")
			err := r.Client.Update(ctx, createdEnvoyFilter, &client.UpdateOptions{})
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{RequeueAfter: 60 * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LocalRateLimitConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ratelimitv1alpha1.LocalRateLimitConfig{}).
		Complete(r)
}
