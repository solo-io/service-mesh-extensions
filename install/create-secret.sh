#!/usr/bin/env bash
GITHUB_TOKEN=${GITHUB_TOKEN:-1234}
NAMESPACE=${NAMESPACE:-sm-marketplace}

kubectl -n $NAMESPACE create secret generic github-token --from-literal=token=$GITHUB_TOKEN
