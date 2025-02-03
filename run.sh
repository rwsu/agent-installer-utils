./hack/build.sh
#RELEASE_IMAGE=quay.io/openshift-release-dev/ocp-release:4.12.2-x86_64 NODE_ZERO_IP=127.0.0.1 ./bin/agent-tui
RELEASE_IMAGE=quay.io/openshift-release-dev/ocp-release:4.12.2-x86_64 NODE_ZERO_IP= ./bin/agent-tui
#RELEASE_IMAGE=virthost.ostest.test.metalkube.org:5000/localimages/local-release-image:4.12.0-rc.7-x86_64 ./bin/agent-tui 
