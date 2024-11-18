package main

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type MyReconciler struct{}

func (a *MyReconciler) Reconcile(
	ctx context.Context,
	req reconcile.Request,
) (reconcile.Result, error) {
	fmt.Printf("reconcile %v\n", req)
	return reconcile.Result{}, nil

}
