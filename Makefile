all: build-local

build-local-deps:
	go mod download

build-local: build-local-deps
	rm -f http-server
	go build -o http-server

build-docker:
	docker build -t local/http-server .

run-k8s-kind: build-docker
	if [ ! "$$(kind get clusters |grep -x fuzzbuzz)" = "fuzzbuzz" ]; then kind create cluster --name fuzzbuzz --config kind.yaml; fi
	kind load docker-image local/http-server:latest --name fuzzbuzz
	helm repo add ingress-nginx 'https://kubernetes.github.io/ingress-nginx'
	helm repo add bitnami 'https://charts.bitnami.com/bitnami'
	helm repo update
	helm upgrade --install ingress-nginx ingress-nginx/ingress-nginx -f kind-ingress-controller.yaml --create-namespace --namespace ingress-nginx --set controller.setAsDefaultIngress=true --set controller.metrics.enabled=true --set controller.service.httpsPort.enable=false --set controller.service.httpPort.enable=true --set controller.service.httpPort.port=80
	kubectl delete -A ValidatingWebhookConfiguration ingress-nginx-admission
	kubectl create ns http-server
	kubectl config set-context --current --namespace=http-server
	helm install redis bitnami/redis --set auth.enabled=false
	helm install http-server helm/http-server