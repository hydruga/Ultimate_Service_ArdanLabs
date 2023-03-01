SHELL := /bin/bash
# =================================================================
# Testing Running System
#
# For testing a simple query on the system. Don't forget to `make seed` first.
# curl --user "admin@example.com:gophers" http://localhost:3000/v1/users/token
# export TOKEN="COPY TOKEN STRING FROM LAST CALL"
# curl -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/users/1/2

# Access metrics directly (4000) or through the sidecar (3001)
# go install github.com/divan/expvarmon@latest
# Run in terminal below
# expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"
# hey -m GET -c 100 -n 10000 http://localhost:3000/v1/test

# To generate a private/public key PEM file.
# openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# openssl rsa -pubout -in private.pem -out public.pem
# ./sales-admin genkey
#
# Testing Auth
# curl -il http://localhost:3000/v1/testauth
# curl -il -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/testauth
# 
# Database Access
# dblab --host 0.0.0.0 --user postgres --db postgres --pass postgres --ssl disable --port 5432 --driver postgres
#
#++++++++++++++++++++ SETUP DEV FOR MAC ++++++++++++++++++++++++++++
dev.setup.mac:
	brew update
	brew list kind || brew install kind
	brew list kubectl || brew install kubectl
	brew list kustomize || brew install kustomize
	brew list pgcli || brew install pgcli
	brew list hey || brew install hey

git-add:
	git add -A
	git commit -m "New files added"
# ========================================================================

run:
	go run app/services/sales-api/main.go | go run app/tooling/logfmt/main.go

admin: 
	go run app/tooling/admin/main.go


# =========================================================================
# Running tests within the local computer.
# walk through entire project tree " ./... "
# count=1 means "ignore cache"
# staticcheck means "lint things"
test:
	go test ./... -count=1
	staticcheck -checks=all ./...

# =========================================================================
# Building Containers

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
	kustomize build zarf/k8s/kind/database-pod | kubectl apply -f -
	kubectl wait --namespace=database-system --timeout=120s --for=condition=Available deployment/database-pod
	kustomize build zarf/k8s/kind/sales-pod | kubectl apply -f -

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --all-namespaces

# We pipe to our logfmt just for human readable logging here, otherwise we get structured logging.
kind-logs:
	kubectl logs -l app=sales -f --tail=100 | go run app/tooling/logfmt/main.go

kind-restart:
	kubectl rollout restart deployment sales-dep

kind-status-sales:
	kubectl get pods -o wide -w

kind-status-db:
	kubectl get pods -o wide -w --namespace=database-system

kind-update: all kind-load kind-restart

kind-update-apply: all kind-load kind-apply

kind-describe:
	kubectl describe pod -l app=sales

#++++++++++++++++++++ GO MODULE SUPPORT ++++++++++++++++++++++++++++++

tidy:
	go mod tidy
#go mod vendor