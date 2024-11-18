# 撰寫 Operator

## 運行最簡單的Operator

```bash 
cd /workspaces/113-1-ntcu-kubernetes-operator-workshop/03-OPERATOR/operator
go mod tidy
go build
```

## 撰寫符合業務邏輯的Operator

```bash
cd /workspaces/113-1-ntcu-kubernetes-operator-workshop/03-OPERATOR/operator

sed -i '/Complete(&MyReconciler{})/c\
		Owns(&corev1.Service{}).\
		Owns(&appsv1.Deployment{}).\
		Complete(&WebReconciler{\
			client: mgr.GetClient(),\
			scheme: mgr.GetScheme(),\
		})' main.go

sed -i '/webv1 "operator\/pkg\/apis\/myweb\/v1"/a\
	appsv1 "k8s.io\/api\/apps\/v1"\
	corev1 "k8s.io\/api\/core\/v1"' main.go

```


```bash
cd /workspaces/113-1-ntcu-kubernetes-operator-workshop/03-OPERATOR/operator

cat << EOF > go.mod
module operator

go 1.22

require (
	k8s.io/api v0.29.2
	k8s.io/apimachinery v0.29.2
	k8s.io/client-go v0.29.2
	sigs.k8s.io/controller-runtime v0.17.6
)
EOF
```


```bash 
cd /workspaces/113-1-ntcu-kubernetes-operator-workshop/03-OPERATOR/operator
go mod tidy
go build
```