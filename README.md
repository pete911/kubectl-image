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

## Example

```
kubectl-image

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
