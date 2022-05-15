CERTS_DIR=../../../docker/k8s
KUBE_API_PORT=6443
DOMAIN=kube.dev
CLUSTER_NAME=dev-cluster
CLUSTER_USER=cluster-admin

CA_CERT_PATH=${CERTS_DIR}/rootCA.crt
CA_KEY_PATH=${CERTS_DIR}/rootCA.key
CERT_KEY_PATH=${CERTS_DIR}/common_cert_key_for_all.key
CERT_PATH=${CERTS_DIR}/common_cert_for_all.crt

KUBECONFIG=${CERTS_DIR}/kubeconfig-dev

codegen:
	./hack/update-codegen.sh

generate-kubeconfig:
	touch ${KUBECONFIG}
	kubectl config set-credentials ${CLUSTER_USER} \
		--kubeconfig=${KUBECONFIG} \
		--client-certificate=${CERT_PATH} \
		--client-key=${CERT_KEY_PATH} \
		--embed-certs=true
	kubectl config set-cluster ${CLUSTER_NAME} \
		--kubeconfig=${KUBECONFIG} \
		--certificate-authority=${CA_CERT_PATH} \
		--server=https://${DOMAIN}:${KUBE_API_PORT}
	kubectl config set-context ${CLUSTER_NAME} \
		--kubeconfig=${KUBECONFIG} \
		--cluster=${CLUSTER_NAME} \
		--user=${CLUSTER_USER}
	kubectl config use-context ${CLUSTER_NAME} --kubeconfig=${KUBECONFIG}

create-crd:
	kubectl create -f ./crd.yaml --kubeconfig=${KUBECONFIG}
