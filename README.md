# kubectl-image

[![pipeline](https://github.com/pete911/kubectl-image/actions/workflows/pipeline.yml/badge.svg)](https://github.com/pete911/kubectl-image/actions/workflows/pipeline.yml)

List container images running in kubernetes cluster. This tool makes it easy to see what images are running in the
kubernetes cluster. Similar thing can be achieved by:

```
kubectl get pods --all-namespaces -o jsonpath="{.items[*].spec.containers[*].image}" |\
tr -s '[[:space:]]' '\n' |\
sort |\
uniq -c
```

But this tool prints more information and the output is sorted, which makes it easier to see:
- what pods and containers are using the image
- output is sorted by registry, repository, tag (or digest/id) and id (if the container is running):
  - easy to spot if image has been re-tagged (same tag, but multiple IDs)
  - easy to spot if multiple versions of the image are running (outdated)

## Usage

```
Available Commands:
  get         get images
  help        Help about any command
  list        list images
  version     print version

Flags:
  -A, --all-namespaces          all kubernetes namespaces
  -h, --help                    help for list
      --kubeconfig string       path to kubeconfig file (default "~/.kube/config")
      --log-level string        log level - debug, info, warn, error (default "warn")
  -n, --namespace string        kubernetes namespace (default "default")
```

- list images in all namespaces `kubectl image list -A`
- list images in a namespaces `kubectl image list -n kube-system`

## Example

### list images
```
kubectl-image list -n kube-system

REGISTRY                                      REPOSITORY                       TAG                         ID                                                 SIZE      PODS  FAILED  RESTART
602401143452.dkr.ecr.eu-west-2.amazonaws.com  amazon-k8s-cni                   v1.14.1-eksbuild.1          sha256:0aea2419c512f8d41a28d9c4c203247582e434a9..  44.05MB   1     0       5
602401143452.dkr.ecr.eu-west-2.amazonaws.com  amazon-k8s-cni                   v1.14.1-eksbuild.1          sha256:60e1f62f53dc02d5bd1df3be1e32756471205261..  44.05MB   3     0       0
602401143452.dkr.ecr.eu-west-2.amazonaws.com  amazon-k8s-cni-init              v1.14.1-eksbuild.1          sha256:2e23a3ecc3fbb541a474a6096cd5ec7ebf91ec6e..  59.66MB   1     0       5
602401143452.dkr.ecr.eu-west-2.amazonaws.com  amazon-k8s-cni-init              v1.14.1-eksbuild.1          sha256:7f5a193cf10e73fc14121aa2fc2f81361aeb9f6c..  59.66MB   3     0       0
602401143452.dkr.ecr.eu-west-2.amazonaws.com  amazon/aws-network-policy-agent  v1.0.2-eksbuild.1           sha256:71fbb862ba51217f4c8a22502cba6fa8baa098b8..  291.58MB  4     0       5
602401143452.dkr.ecr.eu-west-2.amazonaws.com  eks/coredns                      v1.10.1-eksbuild.2          sha256:b5d40542d3a72b3c709de62c40f37f498beba800..  16.22MB   2     0       0
602401143452.dkr.ecr.eu-west-2.amazonaws.com  eks/kube-proxy                   v1.28.1-minimal-eksbuild.1  sha256:5ab5ddb93ba8833e9090ca7ea4780ade10812890..  30.63MB   1     0       5
602401143452.dkr.ecr.eu-west-2.amazonaws.com  eks/kube-proxy                   v1.28.1-minimal-eksbuild.1  sha256:b582776353eee9ba9638f8a487cf7670f3bdf4ea..  30.63MB   3     0       0
registry.k8s.io                               metrics-server/metrics-server    v0.6.4                      sha256:ee4304963fb035239bb5c5e8c10f2f38ee80efc1..  29.96MB   1     0       1
```

In the above example, we can see we have multiple images with the same tag, but different id e.g. `amazon-k8s-cni`.
We can also see that one of them had multiple containers restarted.

### Get image info

To follow from the example above, we can get more information about the image, here we can see that the `amazon-k8s-cni`
image has changed and the pods that started later use different image. This should help us to debug if the node is on
different architecture, or the image has been re-tagged etc.

```
kubectl-image get -n kube-system

registry:   602401143452.dkr.ecr.eu-west-2.amazonaws.com
  repository: amazon-k8s-cni
    id: sha256:0aea2419c512f8d41a28d9c4c203247582e434a90e77a5bcd5beb22cbe7a0a4a tags: v1.14.1-eksbuild.1
    CONTAINER  RESTART  STATE    MESSAGE                          INIT   POD                         PHASE    NODE                                     NODE CREATED
    aws-node   5        running  started at 2023-11-16T00:02:09Z  false  kube-system/aws-node-znbpk  Running  ip-10-0-0-64.eu-west-2.compute.internal  2023-10-13T13:32:32+01:00
    id: sha256:60e1f62f53dc02d5bd1df3be1e327564712052617b05f691cfe322bd2a152505 tags: v1.14.1-eksbuild.1
    CONTAINER  RESTART  STATE    MESSAGE                          INIT   POD                         PHASE    NODE                                      NODE CREATED
    aws-node   0        running  started at 2023-11-22T09:39:23Z  false  kube-system/aws-node-s7299  Running  ip-10-0-1-144.eu-west-2.compute.internal  2023-11-22T09:39:11Z
    aws-node   0        running  started at 2023-11-22T09:39:25Z  false  kube-system/aws-node-w9v27  Running  ip-10-0-0-59.eu-west-2.compute.internal   2023-11-22T09:39:12Z
    aws-node   0        running  started at 2023-11-22T09:39:23Z  false  kube-system/aws-node-q4xds  Running  ip-10-0-2-183.eu-west-2.compute.internal  2023-11-22T09:39:11Z
  repository: amazon-k8s-cni-init
    id: sha256:2e23a3ecc3fbb541a474a6096cd5ec7ebf91ec6ebe4aac3ddc8ff3c18cd6d242 tags: v1.14.1-eksbuild.1
    CONTAINER         RESTART  STATE       MESSAGE                 INIT  POD                         PHASE    NODE                                     NODE CREATED
    aws-vpc-cni-init  5        terminated  exit code: 0 Completed  true  kube-system/aws-node-znbpk  Running  ip-10-0-0-64.eu-west-2.compute.internal  2023-10-13T13:32:32+01:00
    id: sha256:7f5a193cf10e73fc14121aa2fc2f81361aeb9f6ca7edb30f5be3ee7f5ef47ec8 tags: v1.14.1-eksbuild.1
    CONTAINER         RESTART  STATE       MESSAGE                 INIT  POD                         PHASE    NODE                                      NODE CREATED
    aws-vpc-cni-init  0        terminated  exit code: 0 Completed  true  kube-system/aws-node-q4xds  Running  ip-10-0-2-183.eu-west-2.compute.internal  2023-11-22T09:39:11Z
    aws-vpc-cni-init  0        terminated  exit code: 0 Completed  true  kube-system/aws-node-s7299  Running  ip-10-0-1-144.eu-west-2.compute.internal  2023-11-22T09:39:11Z
    aws-vpc-cni-init  0        terminated  exit code: 0 Completed  true  kube-system/aws-node-w9v27  Running  ip-10-0-0-59.eu-west-2.compute.internal   2023-11-22T09:39:12Z
  ...
```

## download

- [binary](https://github.com/pete911/kubectl-image/releases)

## build/install

### brew

- add tap `brew tap pete911/tap`
- install `brew install kubectl-image`

### go

[go](https://golang.org/dl/) has to be installed.
- build `go build`
- install `go install`

## release

Releases are published when the new tag is created e.g.
`git tag -m "add new feature" v1.0.0 && git push --follow-tags`
