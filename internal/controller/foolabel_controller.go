/*
Copyright 2025.

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

package controller

import (
	"context"
	"time"

	foogroupv1 "github.com/git-tux/k8s-foo-controller/api/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// FooLabelReconciler reconciles a FooLabel object
type FooLabelReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=foogroup.foo.controller,resources=foolabels,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=foogroup.foo.controller,resources=foolabels/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=foogroup.foo.controller,resources=foolabels/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the FooLabel object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.4/pkg/reconcile

func (r *FooLabelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = klog.FromContext(ctx)

	var foolabel foogroupv1.FooLabel
	label := foolabel.Spec.Label
	labelvalue := foolabel.Spec.Value

	if err := r.Get(ctx, req.NamespacedName, &foolabel); err != nil {
		if apierrors.IsNotFound(err) {
			// If FooLabel object was deleted, we have to remove the label from pods
			return ctrl.Result{}, nil
		}
		klog.Error(err, "unable to fetch Foo Label")
		return ctrl.Result{}, err
	}

	var podList corev1.PodList
	if err := r.Client.List(ctx, &podList, client.InNamespace(req.Namespace)); err != nil {
		klog.Error(err, "unable to list pods")
		return ctrl.Result{}, err
	}

	for _, pod := range podList.Items {
		if pod.Labels == nil {
			pod.Labels = make(map[string]string)
			pod.Labels[label] = labelvalue
			klog.Info("Added label on pod", pod.Name)
		} else {
			if pod.Labels[label] != labelvalue {
				pod.Labels[label] = labelvalue
				klog.Info("Added label on pod", pod.Name)
			}
		}
	}

	return ctrl.Result{RequeueAfter: time.Minute * 2}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FooLabelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&foogroupv1.FooLabel{}).
		Named("foolabel").
		Complete(r)
}
