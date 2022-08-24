GOLANGCI_IMAGE=golangci/golangci-lint:v1.51.2-alpine

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

trust-ca-cert:
	sudo -S cp ${CA_CERT_PATH} ${CSR_PATH} ${CERT_PATH} /usr/local/share/ca-certificates/
	sudo update-ca-certificates

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

run-test-pod:
	kubectl run test-busybox --image=busybox:1.35.0 --kubeconfig=./docker/k8s/kubeconfig-dev
remove-test-pod:
	kubectl delete pod test-busybox --kubeconfig ./docker/k8s/kubeconfig-dev

test:
	docker-compose up -d && \
		KUBECONFIG=../../../docker/k8s/kubeconfig-dev \
		KUBE_API_SERVER_URL=https://localhost:6443 \
		go test -race -v -coverprofile=coverage.txt -covermode=atomic ./...

lint:
	docker run --rm -v ${PWD}:/app -w /app ${GOLANGCI_IMAGE} golangci-lint run --fix --timeout 20m --sort-results

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