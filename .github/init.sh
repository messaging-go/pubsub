#!/usr/bin/env bash
set -euo pipefail

TARGET="$1"
TARGET_SNAKE_CASE="${TARGET//-/_}"

# replace the repo names
git ls-files -z -- ':!:.github/*' | xargs -0 sed -i \
    -e "s#messaging-go/integration-template#messaging-go/$TARGET#g" \
    -e "s/integration-template/$TARGET/g" \
    -e "s/integration_template/$TARGET_SNAKE_CASE/g"
git mv .idea/integration-template.iml .idea/$TARGET.iml

go install golang.org/x/tools/cmd/goimports@latest
"$(go env GOPATH)/bin/goimports" -w .

rm -f .github/workflows/init.yaml
