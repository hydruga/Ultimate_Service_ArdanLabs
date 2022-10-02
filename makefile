SHELL := /bin/bash

#++++++++++++++++++++ SETUP DEV FOR MAC ++++++++++++++++++++++++++++
dev.setup.mac:
	brew update
	brew list kind || brew install kind
	brew list kubectl || brew install kubectl
	brew list kustomize || brew install kustomize
	brew list pgcli || brew install pgcli

git-add:
	git add -A
	git commit -m "New files added"

run:
	go run main.go

VERSION := 1.0

all: sales-api

sales-api:
	docker build \
		-f zarf/docker/dockerfile.sales-api \
		-t sales-api-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u + "%Y-%m-%dT%H:%M:%SZ"` \
		.

#=====================================================================
# KIND image release info at project: github.com/kubernetes-sigs/kind/releases/tag/[your version of kind]
# Currently using k8s version for github.com/kubernetes-sigs/kind/releases/tag/v0.16.0
# Running from within k8s/kind

KIND_CLUSTER := ardan-starter-cluster 

kind-up:
	kind create cluster \
		--image kindest/node:v1.24.6@sha256:97e8d00bc37a7598a0b32d1fabd155a96355c49fa0d4d4790aab0f161bf31be1 \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yml
	kubectl config set-context --current --namespace=sales-system 

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

#Use kustomize to replace our VERSION from VERSION in make file here
kind-load:
	cd zarf/k8s/kind/sales-pod; kustomize edit set image sales-api-image=sales-api-amd64:$(VERSION)
	kind load docker-image sales-api-amd64:$(VERSION) --name $(KIND_CLUSTER)

# Can use cat for file then use as input for apply -f
# cat zarf/k8s/base/sales-pod/base-sales.yml | kubectl apply -f -
kind-apply:
	kustomize build zarf/k8s/kind/sales-pod | kubectl apply -f -

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --all-namespaces

kind-logs:
	kubectl logs -l app=sales -f --tail=100 

kind-restart:
	kubectl rollout restart deployment sales-dep

kind-status-sales:
	kubectl get pods -o wide -w 

kind-update: all kind-load kind-restart

kind-update-apply: all kind-load kind-apply

kind-describe:
	kubectl describe pod -l app=sales

#++++++++++++++++++++ GO MODULE SUPPORT ++++++++++++++++++++++++++++++

tidy:
	go mod tidy
	go mod vendor