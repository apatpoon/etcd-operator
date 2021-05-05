package controllers

import (
	"context"
	"fmt"
	"reflect"

	corev1 "k8s.io/api/core/v1"

	etcdv1alpha1 "github.com/Simonpoon93/etcd-operator/api/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Action 定义的执行动作接口
type Action interface {
	Execute(ctx context.Context) error
}

type PatchStatus struct {
	client   client.Client
	original runtime.Object
	new      *etcdv1alpha1.EtcdBackup
}

// Execute PatchStatus实现Action接口
func (s *PatchStatus) Execute(ctx context.Context) error {
	// 判断新旧对象是否同一个 是则直接返回不进行更新
	if reflect.DeepEqual(s.original, s.new) {
		return nil
	}
	// 更新状态 使用新object和旧object进行merge
	if err := s.client.Status().Patch(ctx, s.new, client.MergeFrom(s.original)); err != nil {
		return fmt.Errorf("patching status error: %s", err)
	}
	return nil
}

// CreateObject 创建一个新的资源对象
type CreateObject struct {
	client client.Client
	obj    *corev1.Pod
}

// Execute PatchStatus实现Action接口
func (o *CreateObject) Execute(ctx context.Context) error {
	if err := o.client.Create(ctx, o.obj); err != nil {
		return fmt.Errorf("create object error: %s", err)
	}
	return nil
}
