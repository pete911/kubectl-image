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
kubectl image list

Flags:
  -A, --all-namespaces          all kubernetes namespaces
      --field-selector string   kubernetes field selector
  -h, --help                    help for list
      --kubeconfig string       path to kubeconfig file (default "~/.kube/config")
  -l, --label string            kubernetes label
  -n, --namespace string        kubernetes namespace (default "default")
      --size                    print image size (default true)
```

- get all images in all namespaces `kubectl image list -A`
- get all images in a namespaces `kubectl image list -n kube-system`
- select pod images by pod label `kubectl image list -A -l k8s-app=kube-dns`
- specific pod `kubectl image list -n kube-system kube-dns-66bff467f8-7mz46`

## Example

```
kubectl image list -A

registry: registry.k8s.io
  coredns/coredns
    Tag/ID: v1.9.3      Size: 13.42MB
    ID:     sha256:b19406328e70dd2f6a36d6dbe4e867b0684ced2fdeb2f02ecb54ead39ec0bac0
            [namespace] kube-system [container] coredns [pod] coredns-565d847f94-kfs8m [pod-phase] Running
            [namespace] kube-system [container] coredns [pod] coredns-565d847f94-r9wvh [pod-phase] Running
  etcd
    Tag/ID: 3.5.4-0     Size: 81.12MB
    ID:     sha256:8e041a3b0ba8b5f930b1732f7e2ddb654b1739c89b068ff433008d633a51cd03
            [namespace] kube-system [container] etcd [pod] etcd-kind-control-plane [pod-phase] Running
  kube-apiserver
    Tag/ID: v1.25.3     Size: 74.21MB
    ID:     sha256:c666c2ddbc056f8aba649a2647a26d3f6224bce857613b91be6075c88ca963a1
            [namespace] kube-system [container] kube-apiserver [pod] kube-apiserver-kind-control-plane [pod-phase] Running
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
