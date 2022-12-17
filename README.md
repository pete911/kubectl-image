# kubectl-image
List container images running in kubernetes cluster. This tool makes it easy to see what images are running in
kubernetes cluster. Similar thing can be achieved by:

```
kubectl get pods --all-namespaces -o jsonpath="{.items[*].spec.containers[*].image}" |\
tr -s '[[:space:]]' '\n' |\
sort |\
uniq -c
```

But this tool makes it easier to see what pods and containers are using the image. Output is sorted by registry,
repository and then tag (or digest) which makes it easy to spot if containers are running outdated versions (more than
one tag per image repository).

## Usage

```
kubectl-image list

Flags:
  -A, --all-namespaces          all kubernetes namespaces
      --field-selector string   kubernetes field selector
  -h, --help                    help for list
      --kubeconfig string       path to kubeconfig file (default "~/.kube/config")
  -l, --label string            kubernetes label
  -n, --namespace string        kubernetes namespace (default "default")
```

- get all images in all namespaces `kubectl-image list -A`
- get all images in a namespaces `kubectl-image list -n kube-system`
- select pod images by pod label `kubectl-image list -A -l k8s-app=kube-dns`
- specific pod `kubectl-image list -n kube-system kube-dns-66bff467f8-7mz46`

## Example

```
kubectl-image list -A

registry: registry.k8s.io
  coredns/coredns
    v1.9.3
      [namespace] kube-system [container] coredns [pod] coredns-565d847f94-kfs8m [phase] Running
      [namespace] kube-system [container] coredns [pod] coredns-565d847f94-r9wvh [phase] Running
  etcd
    3.5.4-0
      [namespace] kube-system [container] etcd [pod] etcd-kind-control-plane [phase] Running
registry: -
  kindest/kindnetd
    v20221004-44d545d1
      [namespace] kube-system [container] kindnet-cni [pod] kindnet-dfhv7 [phase] Running
  kindest/local-path-provisioner
    v0.0.22-kind.0
      [namespace] local-path-storage [container] local-path-provisioner [pod] local-path-provisioner-684f458cdd-47sfb [phase] Running
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
