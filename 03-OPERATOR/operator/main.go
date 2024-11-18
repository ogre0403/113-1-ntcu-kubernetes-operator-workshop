package main

import (
	"context"
	webv1 "operator/pkg/apis/myweb/v1"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var scheme = runtime.NewScheme()

func init() {
	log.SetLogger(zap.New())
	clientgoscheme.AddToScheme(scheme)
	webv1.AddToScheme(scheme)
}

func main() {

	mgr, err := manager.New(
		config.GetConfigOrDie(),
		manager.Options{
			Scheme: scheme,
		},
	)

	if err != nil {
		panic(err)
	}

	err = builder.
		ControllerManagedBy(mgr).
		For(&webv1.MyWeb{}).
		Complete(&MyReconciler{})

	if err != nil {
		panic(err)
	}

	err = mgr.Start(context.Background())

	if err != nil {
		panic(err)
	}
}
