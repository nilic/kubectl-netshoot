package cmd

import (
	"strings"
)

const (
	globalFlags = `
Global Flags:
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
    --user string                    The name of the kubeconfig user to use`

	debugFlags = `
Flags:
    --image-tag string               netshoot container image tag to use (default "latest")
-h, --help                           help for debug`

	debugShort = "Debug using an ephemeral container in an existing pod or on a node"

	debugUsage = `
Usage:
    kubectl netshoot debug (POD | TYPE[[.VERSION].GROUP]/NAME) [flags] [ -- COMMAND [args...] ]`

	debugExamples = `
Examples:
    # debug using an ephemeral container in an existing pod
    kubectl netshoot debug my-existing-pod

    # debug with a specific netshoot image
    kubectl netshoot debug my-existing-pod --image-tag v0.5

    # run a command in netshoot container instead of shell
    kubectl netshoot debug mypod -- curl localhost:8443

    # create a debug session on a node
    # netshoot will run in the host namespaces and have host's filesystem mounted at /host
    kubectl netshoot debug node/my-node`

	runFlags = `
Flags:
    --host-network                   (applicable to "run" command only) whether to spin up netshoot on the host's network namespace
    --image-tag string               netshoot container image tag to use (default "latest")
-h, --help                           help for run`

	runShort = "Run a throwaway pod for troubleshooting"

	runUsage = `
Usage:
    kubectl netshoot run NAME [flags] [ -- COMMAND [args...] ]`

	runExamples = `
Examples:
    # spin up a throwaway pod for troubleshooting
    kubectl netshoot run tmp-shell

    # spin up a throwaway pod with a specific netshoot image
    kubectl netshoot run tmp-shell --image-tag v0.5

    # run a command in netshoot container instead of shell
    kubectl netshoot run tmp-shell -- ping 8.8.8.8

    # spin up a netshoot pod on the host's network namespace
    kubectl netshoot run tmp-shell --host-network`
)

var (
	debugHelp = strings.Join([]string{debugShort, debugExamples, debugUsage, debugFlags, globalFlags}, "\n")
	runHelp   = strings.Join([]string{runShort, runExamples, runUsage, runFlags, globalFlags}, "\n")
)
