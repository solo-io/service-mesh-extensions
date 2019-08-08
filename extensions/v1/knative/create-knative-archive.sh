#!/usr/bin/env bash

# get vanilla manifest
wget https://github.com/knative/serving/releases/download/v0.8.0/serving-core.yaml
tar zcvf serving-core-0.8.0.tar.gz serving-core.yaml
rm serving-core.yaml

# get istio manifest
wget https://github.com/knative/serving/releases/download/v0.8.0/serving.yaml
tar zcvf serving-0.8.0.tar.gz serving.yaml
rm serving.yaml
