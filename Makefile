GOLANGCI_IMAGE=golangci/golangci-lint:v2.0.2-alpine

HOST_IP=`hostname -I | awk '{print $$1}'`

K8S_CLUSTER_NAME=dev-cluster
ASSETS_DIR=/home/erema/h/go_sandbox/k8s/dev_environment/assets
CLUSTER_NAME=dev-cluster
CLUSTER_USER=cluster-admin
KUBE_API_PORT=6443
DOMAIN=localhost
CA_CERT_PATH=${ASSETS_DIR}/rootCA.crt
CA_KEY_PATH=${ASSETS_DIR}/rootCA.key
CSR_PATH=${ASSETS_DIR}/csr_for_common_cert_for_all.csr
CERT_KEY_PATH=${ASSETS_DIR}/common_cert_key_for_all.key
CERT_PATH=${ASSETS_DIR}/common_cert_for_all.crt

KUBECONFIG_PATH=${ASSETS_DIR}/kubeconfig-dev

gen-sa-certs:
	openssl genrsa -out ${CA_KEY_PATH} 2048
	openssl req -x509 -new -nodes -key ${CA_KEY_PATH} -days 10000 -out ${CA_CERT_PATH} -config csr.conf
	openssl genrsa -out ${CERT_KEY_PATH} 2048
	openssl req -new -key ${CERT_KEY_PATH} -out ${CSR_PATH} -config csr.conf
	openssl x509 -req -in ${CSR_PATH} -CA ${CA_CERT_PATH} -CAkey ${CA_KEY_PATH} \
		-CAcreateserial -out ${CERT_PATH} -days 10000 \
		-extensions v3_ext -extfile csr.conf -sha256
	openssl req  -noout -text -in ${CSR_PATH}
	openssl x509  -noout -text -in ${CERT_PATH}

trust-ca-cert:
	sudo -S cp ${CA_CERT_PATH} ${CERT_PATH} /usr/local/share/ca-certificates/
	sudo update-ca-certificates -f

gen-kubeconfig:
	touch ${KUBECONFIG_PATH}
	kubectl config set-credentials ${CLUSTER_USER} \
		--kubeconfig=${KUBECONFIG_PATH} \
		--client-certificate=${CERT_PATH} \
		--client-key=${CERT_KEY_PATH} \
		--embed-certs=true
	kubectl config set-cluster ${K8S_CLUSTER_NAME} \
		--kubeconfig=${KUBECONFIG_PATH} \
		--certificate-authority=${CA_CERT_PATH} \
		--server=https://${DOMAIN}:${KUBE_API_PORT}
	kubectl config set-context admin-context \
		--kubeconfig=${KUBECONFIG_PATH} \
		--cluster=${K8S_CLUSTER_NAME} \
		--user=${CLUSTER_USER}
	kubectl config use-context admin-context --kubeconfig=${KUBECONFIG_PATH}

run-test-pod:
	kubectl run test-busybox --image=busybox:1.35.0 --kubeconfig=./docker/k8s/kubeconfig-dev
remove-test-pod:
	kubectl delete pod test-busybox --kubeconfig ./docker/k8s/kubeconfig-dev

test:
	docker compose up -d && \
		KUBECONFIG=../../../docker/k8s/kubeconfig-dev \
		KUBE_API_SERVER_URL=https://localhost:6443 \
		go test -race -v -coverprofile=coverage.txt -covermode=atomic ./...

lint:
	docker run --rm -v ${PWD}:/app -w /app ${GOLANGCI_IMAGE} golangci-lint run --fix --timeout 20m

CNI_PLUGIN_ARCHIVE=cni-plugins-linux-amd64-v1.1.1.tgz
setup-containerd:
	wget https://github.com/containernetworking/plugins/releases/download/v1.1.1/${CNI_PLUGIN_ARCHIVE}
	sudo -S mkdir -p /opt/cni/bin
	sudo tar Cxzvf /opt/cni/bin cni-plugins-linux-amd64-v1.1.1.tgz
	rm ${CNI_PLUGIN_ARCHIVE}
	sudo cp 100-crio-bridge.conf 200-loopback.conf /etc/cni/net.d
	sudo -S bash -c "containerd config default > /etc/containerd/config.toml"
	sudo systemctl restart containerd

running-containers-loop:
	sudo -S bash -c "while sleep 1; do date; ctr --namespace k8s.io containers list; done;"

etcd:
	docker run --rm -p 2379:2379 -p 2380:2380 --name etcd quay.io/coreos/etcd:v3.5.1 /usr/local/bin/etcd \
		--name node1 \
		--initial-advertise-peer-urls http://${HOST_IP}:2380 \
		--listen-peer-urls http://0.0.0.0:2380 \
		--advertise-client-urls http://${HOST_IP}:2379 \
		--listen-client-urls http://0.0.0.0:2379 \
		--initial-cluster node1=http://${HOST_IP}:2380 \
		--log-level debug

K8S_API_SERVER_DEBUG_PORT=62001
K8S_CONTROLLER_MANAGER_DEBUG_PORT=62002
K8S_SCHEDULER_DEBUG_PORT=62003
K8S_KUBELET_DEBUG_PORT=62004

set-oidc-in-kubeconfig:
	curl -d 'client_id=${OIDC_CLIENT_NAME}' \
		-d 'username=admin' \
		-d 'password=admin' \
		-d 'grant_type=password' \
		-d 'client_secret=${OIDC_CLIENT_SECRET}' \
		-d 'scope=openid'  \
		https://localhost:8443/realms/master/protocol/openid-connect/token > tmp_token.json
	kubectl config set-credentials developer-user \
		--auth-provider=oidc \
		--auth-provider-arg=idp-issuer-url="https://localhost:8443/realms/master" \
		--auth-provider-arg=client-id=${OIDC_CLIENT_NAME} \
		--auth-provider-arg=client-secret=${CLIENT_SECRET} \
		--auth-provider-arg=refresh-token=`cat token.json | jq -r .refresh_token` \
		--auth-provider-arg=id-token=`cat token.json | jq -r .id_token` \
		--kubeconfig=${KUBECONFIG_PATH}
	rm tmp_token.json
	kubectl config set-context developer-context \
		--kubeconfig=${KUBECONFIG_PATH} \
		--cluster=${K8S_CLUSTER_NAME} \
		--user=developer-user

# https://about.gitlab.com/blog/2018/06/07/keeping-git-commit-history-clean/
start-changing-git-commit:
	# 1. Go to the previous commit before target commit
	git rebase -i `git log --pretty=%P -n 1 ${TARGET_COMMIT_TO_CHANGE}`
	# 2. Change "pick -> edit" desired commit(first in the list), example:
	# pick 74748f9 CI adding                edit 74748f9 CI adding
	# pick 63f7877 Brunch Sums Problem  =>  pick 63f7877 Brunch Sums Problem
	# ...                                   ...
	# 3. Make needed changes and add to commit changed files, example : git add .github/workflows/lint.yml
	# 4. Run `make finish-changing-git-commit`
finish-changing-git-commit:
	git rebase --continue
	git push --force-with-lease origin master
