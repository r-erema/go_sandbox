GOLANGCI_IMAGE=golangci/golangci-lint:latest-alpine

DATA_DIR=./docker/k8s
CLUSTER_NAME=dev-cluster
CLUSTER_USER=cluster-admin
KUBE_API_PORT=6443
DOMAIN=localhost
CA_CERT_PATH=${DATA_DIR}/rootCA.crt
CA_KEY_PATH=${DATA_DIR}/rootCA.key
CSR_PATH=${DATA_DIR}/cluster.csr
CERT_KEY_PATH=${DATA_DIR}/common_cert_key_for_all.key
CERT_PATH=${DATA_DIR}/common_cert_for_all.crt

KUBECONFIG_PATH=${DATA_DIR}/kubeconfig-dev

gen-sa-certs:
	mkdir -p ${DATA_DIR}
	openssl req -new -newkey rsa:4096 -days 365 -nodes -x509 \
		-subj "/C=BY/ST=Minsk Region/L=Minsk/O=${DOMAIN} Office/CN=${DOMAIN}/subjectAltName=DNS.1=${DOMAIN}" \
		-keyout ${CA_KEY_PATH} -out ${CA_CERT_PATH}
	echo "CA Key ${CA_KEY_PATH} is ready"
	echo "CA Cert ${CA_CERT_PATH} is ready"
	openssl genrsa -out ${CERT_KEY_PATH} 2048
	echo "Cert Key ${CERT_KEY_PATH} is ready"
	openssl req -new -key ${CERT_KEY_PATH} \
		-subj "/C=BY/ST=Minsk Region/L=Minsk/O=${DOMAIN} Office/CN=${DOMAIN}/subjectAltName=DNS.1=${DOMAIN}" \
		-out ${CSR_PATH}
	echo "CSR ${CSR_PATH} is ready"
	printf "subjectAltName=DNS:${DOMAIN}" > tmp-ext-file
	openssl x509 -req -extfile tmp-ext-file -in ${CSR_PATH} -days 365 \
		-CA ${CA_CERT_PATH} \
		-CAkey ${CA_KEY_PATH} \
		-CAcreateserial \
		-out ${CERT_PATH}
	echo "cert ${CERT_PATH} is ready"
	rm tmp-ext-file

gen-kubeconfig:
	touch ${KUBECONFIG_PATH}
	kubectl config set-credentials ${CLUSTER_USER} \
		--kubeconfig=${KUBECONFIG_PATH} \
		--client-certificate=${CERT_PATH} \
		--client-key=${CERT_KEY_PATH} \
		--embed-certs=true
	kubectl config set-cluster ${CLUSTER_NAME} \
		--kubeconfig=${KUBECONFIG_PATH} \
		--certificate-authority=${CA_CERT_PATH} \
		--server=https://${DOMAIN}:${KUBE_API_PORT}
	kubectl config set-context ${CLUSTER_NAME} \
		--kubeconfig=${KUBECONFIG_PATH} \
		--cluster=${CLUSTER_NAME} \
		--user=${CLUSTER_USER}
	kubectl config use-context ${CLUSTER_NAME} --kubeconfig=${KUBECONFIG_PATH}

test:
	docker-compose up -d && docker-compose exec \
										-e MYSQL_HOST=mysql \
										-e POSTGRES_HOST=postgres \
										-e MONGODB_HOST=mongodb \
										golang go test -race -v ./...

lint:
	docker run --rm -v ${PWD}:/app -w /app ${GOLANGCI_IMAGE} golangci-lint run --fix --timeout 20m --sort-results
