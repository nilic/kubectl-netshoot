# kubectl netshoot 🌠
A [kubectl plugin](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/) to easily spin up and access a [netshoot](https://github.com/nicolaka/netshoot) container. `netshoot` is a network troubleshooting Swiss-army knife which allows you to perform Kubernetes troubleshooting without installing any new packages in your containers or cluster nodes.

## Installation

### krew

Using [krew](https://krew.sigs.k8s.io/) kubectl plugin manager:

```
# add kubectl-netshoot plugin index
kubectl krew index add netshoot https://github.com/nilic/kubectl-netshoot.git

# install netshoot plugin
kubectl krew install netshoot/netshoot
```

### Without krew

Just download the binary for your OS and architecture from the [Releases](https://github.com/nilic/kubectl-netshoot/releases) page and place it in your `PATH`.

## Usage

```
Usage:
  kubectl netshoot [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  debug       Debug using an ephemeral container in an existing pod or on a node
  help        Help about any command
  run         Run a throwaway pod for troubleshooting
  version     Print kubectl-netshoot version

Flags:
  -h, --help                           help for kubectl-netshoot
      --host-network                   ("run" command only) spin up netshoot on the node's network namespace
      --image-name string              netshoot container image to use (default "nicolaka/netshoot")
      --image-tag string               netshoot container image tag to use (default "latest")
```

In addition to these flags, the following `kubectl` flags are available for all `kubectl netshoot` commands:

```
Flags:
      --as string                      Username to impersonate for the operation. User could be a regular user or a service account in a namespace.
      --as-group stringArray           Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --as-uid string                  UID to impersonate for the operation.
      --cache-dir string               Default cache directory (default "$HOME/.kube/cache")
      --certificate-authority string   Path to a cert file for the certificate authority
      --client-certificate string      Path to a client certificate file for TLS
      --client-key string              Path to a client key file for TLS
      --cluster string                 The name of the kubeconfig cluster to use
      --context string                 The name of the kubeconfig context to use
      --disable-compression            If true, opt-out of response compression for all requests to the server
      --insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string              Path to the kubeconfig file to use for CLI requests.
      --match-server-version           Require server version to match client version
  -n, --namespace string               If present, the namespace scope for this CLI request
      --request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
  -s, --server string                  The address and port of the Kubernetes API server
      --tls-server-name string         Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used
      --token string                   Bearer token for authentication to the API server
      --user string                    The name of the kubeconfig user to use
```

## Examples

### Interactive session

Common use-case is to spin up a `netshoot` container with an interactive session and attach to it so you can perform troubleshooting using tools available in `netshoot`.

`run` command is equivalent to `kubectl run --rm` meaning that the pod will be deleted after you exit the session.

```
# spin up a throwaway pod for troubleshooting
kubectl netshoot run tmp-shell

# spin up a throwaway pod with a specific netshoot image
kubectl netshoot run tmp-shell --image-tag v0.5

# spin up a throwaway pod on the node's network namespace
kubectl netshoot run tmp-shell --host-network
```

`debug` command spins up netshoot as an [ephemeral container](https://kubernetes.io/docs/concepts/workloads/pods/ephemeral-containers/) in an existing pod or on a node. Ephemeral container terminates after the interactive session is exited.

```
# debug using an ephemeral container in an existing pod
kubectl netshoot debug my-existing-pod

# debug with a specific netshoot image
kubectl netshoot debug my-existing-pod --image-tag v0.5

# create a debug session on a node
kubectl netshoot debug node/my-node
```

When running a debug session on a node, `netshoot` will run in the node's network namespace and have node's filesystem mounted at `/host`.

### Running one-time commands

Instead of attaching to the shell, you can also run a one-time command directly on the `netshoot` container. The command you want to run is specified after `--`: 

```
kubectl netshoot run tmp-shell -- ping 8.8.8.8
```

```
kubectl netshoot debug mypod -- curl localhost:8443
```

### Troubleshooting with a custom container

By default, plugin spins up `nicolaka/netshoot:latest` container image, but flags `--image-name` and `--image-tag` allow for running an arbitrary container, eg.

```
kubectl netshoot run tmp-shell --image-name busybox --image-tag 1.36.0
```