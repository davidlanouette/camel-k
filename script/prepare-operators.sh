#!/bin/bash

# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

if [ "$#" -lt 2 ]; then
    echo "usage: $0 prepare-operators <release-version> <github-user>"
    exit 1
fi

location=$(dirname $0)
version=$1
gh_user=$2

cd bundle/

mkdir -p k8s-operatorhub/$1/manifests/
mkdir -p k8s-operatorhub/$1/metadata/
mkdir -p k8s-operatorhub/$1/tests/scorecard/
mkdir -p openshift-ecosystem/$1/manifests/
mkdir -p openshift-ecosystem/$1/metadata/
mkdir -p openshift-ecosystem/$1/tests/scorecard/

cp ./manifests/camel.apache.org_builds.yaml k8s-operatorhub/$1/manifests/builds.camel.apache.org.crd.yaml
cp ./manifests/camel.apache.org_camelcatalogs.yaml k8s-operatorhub/$1/manifests/camelcatalogs.camel.apache.org.crd.yaml
cp ./manifests/camel.apache.org_integrationkits.yaml k8s-operatorhub/$1/manifests/integrationkits.camel.apache.org.crd.yaml
cp ./manifests/camel.apache.org_integrationplatforms.yaml k8s-operatorhub/$1/manifests/integrationplatforms.camel.apache.org.crd.yaml
cp ./manifests/camel.apache.org_integrationprofiles.yaml k8s-operatorhub/$1/manifests/integrationprofiles.camel.apache.org.crd.yaml
cp ./manifests/camel.apache.org_integrations.yaml k8s-operatorhub/$1/manifests/integrations.camel.apache.org.crd.yaml
cp ./manifests/camel.apache.org_kamelets.yaml k8s-operatorhub/$1/manifests/kamelets.camel.apache.org.crd.yaml
cp ./manifests/camel.apache.org_pipes.yaml k8s-operatorhub/$1/manifests/pipes.camel.apache.org.crd.yaml
cp ./manifests/camel-k.clusterserviceversion.yaml k8s-operatorhub/$1/manifests/camel-k.v$1.clusterserviceversion.yaml
cp ./metadata/annotations.yaml k8s-operatorhub/$1/metadata/annotations.yaml
cp ./tests/scorecard/config.yaml k8s-operatorhub/$1/tests/scorecard/config.yaml

cp ./manifests/camel.apache.org_builds.yaml openshift-ecosystem/$1/manifests/builds.camel.apache.org.crd.yaml
cp ./manifests/camel.apache.org_camelcatalogs.yaml openshift-ecosystem/$1/manifests/camelcatalogs.camel.apache.org.crd.yaml
cp ./manifests/camel.apache.org_integrationkits.yaml openshift-ecosystem/$1/manifests/integrationkits.camel.apache.org.crd.yaml
cp ./manifests/camel.apache.org_integrationplatforms.yaml openshift-ecosystem/$1/manifests/integrationplatforms.camel.apache.org.crd.yaml
cp ./manifests/camel.apache.org_integrationprofiles.yaml openshift-ecosystem/$1/manifests/integrationprofiles.camel.apache.org.crd.yaml
cp ./manifests/camel.apache.org_integrations.yaml openshift-ecosystem/$1/manifests/integrations.camel.apache.org.crd.yaml
cp ./manifests/camel.apache.org_kamelets.yaml openshift-ecosystem/$1/manifests/kamelets.camel.apache.org.crd.yaml
cp ./manifests/camel.apache.org_pipes.yaml openshift-ecosystem/$1/manifests/pipes.camel.apache.org.crd.yaml
cp ./manifests/camel-k.clusterserviceversion.yaml openshift-ecosystem/$1/manifests/camel-k.v$1.clusterserviceversion.yaml
cp ./metadata/annotations.yaml openshift-ecosystem/$1/metadata/annotations.yaml
cp ./tests/scorecard/config.yaml openshift-ecosystem/$1/tests/scorecard/config.yaml

# Starting sed to replace operator

sed -i 's/camel-k.v/camel-k-operator.v/g' k8s-operatorhub/$1/manifests/camel-k.v$1.clusterserviceversion.yaml
sed -i 's/camel-k.v/camel-k-operator.v/g' openshift-ecosystem/$1/manifests/camel-k.v$1.clusterserviceversion.yaml

# Clone projects
git clone https://github.com/$gh_user/community-operators.git /tmp/operators/community-operators
cp -r k8s-operatorhub/$version /tmp/operators/community-operators/operators/camel-k/.
git clone https://github.com/$gh_user/community-operators-prod.git /tmp/operators/community-operators-prod
cp -r openshift-ecosystem/$version /tmp/operators/community-operators-prod/operators/camel-k/.

# Community operators
cd /tmp/operators/community-operators
git checkout -b feat/v$version
git add operators/camel-k/$version
git commit -s -m "operator camel-k ($version)"
git remote add upstream https://github.com/k8s-operatorhub/community-operators -f
git pull --rebase upstream main
git push --set-upstream origin feat/v$version

# Community operators PROD
cd /tmp/operators/community-operators-prod
git checkout -b feat/v$version
git add operators/camel-k/$version
git commit -s -m "operator camel-k ($version)"
git remote add upstream https://github.com/redhat-openshift-ecosystem/community-operators-prod -f
git pull --rebase upstream main
git push --set-upstream origin feat/v$version

echo "### You need to create PRs manually:"
echo "--> https://github.com/$gh_user/community-operators/pull/new/feat/v$version"
echo "--> https://github.com/$gh_user/community-operators-prod/pull/new/feat/v$version"
