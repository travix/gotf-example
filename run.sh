#!/usr/bin/env bash
set -eo pipefail

if [[ -z "${CI}" ]]; then
  trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM EXIT
fi

. gen.sh

echo "[run.sh] Building and installing example provider"
pushd provider
go build -o terraform-provider-example
os_arch=$(go env GOOS)_$(go env GOARCH)
mkdir -p "$HOME/.terraform.d/plugins/travix.com/providers/example/1.0.0/${os_arch}/"
mv terraform-provider-example "$HOME/.terraform.d/plugins/travix.com/providers/example/1.0.0/${os_arch}/"
popd

export TF_VAR_key_id="key-one"
TF_VAR_secret_key="$(head /dev/urandom | LC_ALL=C tr -dc A-Za-z0-9 | head -c24)"
export TF_VAR_secret_key
# setting these variables as they are shared between the provider and the example server for the sake of simplicity

pushd example-server
go run . &
server_pid=$!
sleep 1
echo "[run.sh] Staring grpc server pid: $server_pid"
popd

echo "[run.sh] Running terraform"
pushd tfscript
#export TF_LOG_PROVIDER=INFO
rm -rf .terraform .terraform.lock.hcl plan.out terraform.tfstate terraform.tfstate.backup
echo "[run.sh] terraform init"
terraform init
echo "[run.sh] terraform valid"
terraform validate
echo "[run.sh] terraform plan"
terraform plan -out=plan.out
echo "[run.sh] terraform apply"
terraform apply plan.out
popd
