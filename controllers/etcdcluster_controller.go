/*
Copyright 2021 Simonpoon93.

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
	"k8s.io/client-go/util/retry"

	appsv1 "k8s.io/api/apps/v1"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	corev1 "k8s.io/api/core/v1"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	etcdv1alpha1 "github.com/Simonpoon93/etcd-operator/api/v1alpha1"
)

// EtcdClusterReconciler reconciles a EtcdCluster object
type EtcdClusterReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=etcd.oschina.cn,resources=etcdclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=etcd.oschina.cn,resources=etcdclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=etcd.oschina.cn,resources=etcdclusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the EtcdCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *EtcdClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := r.Log.WithValues("etcdcluster", req.NamespacedName)

	// your logic here

	// 获取etcdCluster实例
	var etcdCluster etcdv1alpha1.EtcdCluster
	// 使用r.Get访问本地indexer缓存, 而不是直接访问apiserver
	if err := r.Get(ctx, req.NamespacedName, &etcdCluster); err != nil {
		// 返回NotFound的错误 使当前req重新入队
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 已经获取到etcdCluster实例
	// 创建对应的headless svc 和 statefulset
	var svc corev1.Service
	svc.Name = etcdCluster.Name
	svc.Namespace = etcdCluster.Namespace

	// 尝试进行冲突重试
	if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// CreateOrUpdate
		// 调谐: 观察当前与期望状态
		opResult, err := ctrl.CreateOrUpdate(ctx, r.Client, &svc, func() error {
			// 调谐函数必须在这里实现 实际上是去拼装service
			MutateHeadlessSvc(&etcdCluster, &svc)
			// 为svc资源添加ControllerReference 为自定义crd与对应操作资源绑定从属关系
			return controllerutil.SetControllerReference(&etcdCluster, &svc, r.Scheme)
		})
		log.Info("CreateOrUpdate Result", "Service", opResult)
		return err
	}); err != nil {
		return ctrl.Result{}, err
	}

	// 调谐: 观察当前与期望状态
	var sts appsv1.StatefulSet
	sts.Name = etcdCluster.Name
	sts.Namespace = etcdCluster.Namespace

	if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		opResult, err := ctrl.CreateOrUpdate(ctx, r.Client, &sts, func() error {
			// 调谐函数必须在这里实现 实际上是去拼装service
			MutateStatefulSet(&etcdCluster, &sts)
			// 为svc资源添加ControllerReference 为自定义crd与对应操作资源绑定从属关系
			return controllerutil.SetControllerReference(&etcdCluster, &sts, r.Scheme)
		})
		log.Info("CreateOrUpdate Result", "StatefulSet", opResult)
		return err
	}); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EtcdClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&etcdv1alpha1.EtcdCluster{}).
		// 新增对crd操作到的资源进行watch以维护状态
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
