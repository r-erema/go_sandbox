#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
CODEGEN_PKG=${CODEGEN_PKG:-$(ls -d -1 "${GOPATH}"/pkg/mod/k8s.io/code-generator@v0.23.5 2>/dev/null || echo ../code-generator)}
MODULE=github.com/r-erema/go_sendbox

bash "${CODEGEN_PKG}"/generate-groups.sh "deepcopy,client,informer,lister" \
  ${MODULE}/learning/k8s/example1/pkg/generated \
  ${MODULE}/learning/k8s/example1/pkg/apis \
  samplecontroller:v1alpha1 \
  --go-header-file "${SCRIPT_ROOT}"/hack/boilerplate.go.txt

cp -r "${GOPATH}/src/${MODULE}"/learning/k8s/example1/pkg/apis/samplecontroller/v1alpha1 \
      "$(dirname "${BASH_SOURCE[0]}")"/../pkg/apis/samplecontroller

cp -r "${GOPATH}/src/${MODULE}"/learning/k8s/example1/pkg/generated \
      "$(dirname "${BASH_SOURCE[0]}")"/../pkg
