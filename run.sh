#!/usr/bin/env bash
set -euo pipefail
trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM EXIT

./gen.sh

echo "[run.sh] Building and installing example provider"
go build -o terraform-provider-example ./provider/*
os_arch=$(go env GOOS)_$(go env GOARCH)
mkdir -p "$HOME/.terraform.d/plugins/travix.com/providers/example/1.0.0/${os_arch}/"
mv terraform-provider-example "$HOME/.terraform.d/plugins/travix.com/providers/example/1.0.0/${os_arch}/"

go run example-server/* &
server_pid=$!
sleep 1
echo "[run.sh] Staring grpc server pid: $server_pid"

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
