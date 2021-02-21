# Install all Knative components

Source: https://github.com/csantanapr/knative-kind

## Edit cm config-deployment and add `kind.local`:

```
registriesSkippingTagResolving: kind.local,ko.local,dev.local
```

# Knative Serving

Show differences between K8s deployment and knative service (PHP example) YAMLs

## Show autoscaling with hello service installed during knative setup

```
curl http://hello.default.127.0.0.1.nip.io
```

## Show how `kn` can help

```
kn service list
```

## Deploy `hello-go` with `ko`

Explain what it does, then:

```
# export variables
export KIND_CLUSTER_NAME=knative
export KO_DOCKER_REPO=kind.local

# show parsed ko YAML
ko resolve -Pf config/hello-kservice.yaml

# resolve, publish and deploy
ko resolve -Pf config/hello-kservice.yaml| k apply -f -

# show logs
stern -A 'hello-go' -c user-container

# curl 
curl http://hello-go.default.127.0.0.1.nip.io
curl http://hello-go.default.127.0.0.1.nip.io\?name\=michael
```

## Open Octant and inspect what happened

```
octant
```

## Show Traffic Management with `Revisions`

```
while true; do curl http://hello-go.default.127.0.0.1.nip.io ; echo ""; sleep .1; done

# note the exact revision
kubectl get rev

# edit KService by replacing traffic section with:
  traffic:
  - percent: 100
    revisionName: hello-go-r9qn5

# make code changes and deploy new version (revision)
ko resolve -Pf config/hello-kservice.yaml| k apply -f -

# update KService with 90/50 split
  traffic:
  - latestRevision: false
    percent: 90
    revisionName: hello-go-r9qn5
  - latestRevision: false
    percent: 10
    revisionName: hello-go-nfgmv

# if all is well, switch to v2 - note that v1 will be scaled to zero after some time
  traffic:
  - latestRevision: true
    percent: 100
```

# Knative Eventing

## Install Sources for Tanzu

```
export KIND_CLUSTER_NAME=knative
export KO_DOCKER_REPO=kind.local

ko apply -f config
```

## Authenticate with vCenter

Make sure `kn-vsphere` plugin is installed in `~/.config/kn/plugins`

```
kn plugin list
```

Login in to vCenter

```
kn vsphere login --secret-name vsphere-credentials --username administrator@vsphere.local --password 'Admin!23'
```

## Create a `VSphereSource`

```
kn vsphere source --name my-vc --address https://10.161.153.226 --skip-tls-verify --secret-ref vsphere-credentials --sink-uri http://broker-ingress.knative-eventing.svc.cluster.local/default/default

kubectl get vspheresources.sources.tanzu.vmware.com
```

Check if the source deployment is healthy

```
stern -A 'my-vc'
```

## Create an Event Subscriber Function

```
kn service create kn-echo --port 8080 --image embano1/kn-echo:latest --scale 1
```

## Create the Trigger

```
kn trigger create veba-kn-echo --sink ksvc:kn-echo
```

## Stream logs from the function and then trigger an event

```
stern -A 'kn-echo'
```

## Create a `VSphereBinding`

```
kubectl apply -f config/binding.yaml
```

## Create a Deployment to perform pre-authenticated actions against vCenter

```
kubectl apply -f config/govc-deployment.yaml
```

Run a `$SHELL` in the pod, using the injected environment variables from the `VSphereBinding` above

```
kubectl exec -it govc-658b9c6bc8-hmx44 -- /bin/sh -c 'GOVC_URL=$VC_URL GOVC_INSECURE=$VC_INSECURE GOVC_PASSWORD=$VC_PASSWORD GOVC_USERNAME=$VC_USERNAME bash'
```

Execute some commands in the `$SHELL`

```
govc ls /vcqaDC/vm
govc vm.info /vcqaDC/vm/test-01

# this should also trigger an event if kn-echo is still running
govc vm.power --on /vcqaDC/vm/test-01
```