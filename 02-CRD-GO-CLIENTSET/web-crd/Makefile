CRD_NAME := myweb
VERSION := v1
BASE_PATH := web-crd


install-client-gen:
	go install k8s.io/code-generator/cmd/client-gen@v0.29.2

install-deepcopy-gen:
	go install k8s.io/code-generator/cmd/deepcopy-gen@v0.29.2

install-register-gen:
	go install k8s.io/code-generator/cmd/register-gen@v0.29.2

install-informer-gen:
	go install k8s.io/code-generator/cmd/informer-gen@v0.29.2

install-lister-gen:
	go install k8s.io/code-generator/cmd/lister-gen@v0.29.2


install-generator: install-client-gen install-deepcopy-gen install-register-gen install-informer-gen install-lister-gen

generate-deepcopy: install-deepcopy-gen
	deepcopy-gen \
	--input-dirs $(BASE_PATH)/pkg/apis/$(CRD_NAME)/$(VERSION) \
	-O zz_generated.deepcopy \
	--output-base .. \
	--go-header-file \
	./hack/boilerplate.go.txt

generate-clientset: install-client-gen
	client-gen \
	--clientset-name clientset \
	--input-base "" \
	--input $(BASE_PATH)/pkg/apis/$(CRD_NAME)/$(VERSION) \
	--output-package $(BASE_PATH)/pkg/ \
	--output-base .. \
	--go-header-file ./hack/boilerplate.go.txt


generate-register: install-register-gen
	register-gen \
    -O zz_generated.register \
    --go-header-file ./hack/boilerplate.go.txt \
	--input-dirs ${BASE_PATH}/pkg/apis/${CRD_NAME}/${VERSION} \
    --output-base ..

generate-informer: install-informer-gen
	informer-gen \
	--input-dirs $(BASE_PATH)/pkg/apis/$(CRD_NAME)/$(VERSION) \
	--versioned-clientset-package $(BASE_PATH)/pkg/clientset \
	--listers-package $(BASE_PATH)/pkg/listers \
	--output-package $(BASE_PATH)/pkg/informers \
	--output-base .. \
	--go-header-file ./hack/boilerplate.go.txt

generate-lister: install-lister-gen
	lister-gen \
	--input-dirs $(BASE_PATH)/pkg/apis/$(CRD_NAME)/$(VERSION) \
	--output-package $(BASE_PATH)/pkg/listers \
	--output-base .. \
	--go-header-file ./hack/boilerplate.go.txt

auto-generate: generate-deepcopy generate-register generate-clientset generate-informer generate-lister


remove-generated:
	rm -rf pkg/clientset
	rm -rf pkg/informers
	rm -rf pkg/listers
	rm -rf pkg/apis/$(CRD_NAME)/$(VERSION)/zz_generated.deepcopy.go
	rm -rf pkg/apis/$(CRD_NAME)/$(VERSION)/zz_generated.register.go

clean: remove-generated
	rm -rf pkg/
	rm -rf hack/
	rm main.go