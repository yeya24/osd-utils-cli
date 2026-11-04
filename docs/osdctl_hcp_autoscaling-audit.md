## osdctl hcp autoscaling-audit

Audit all HCP management clusters for autoscaling readiness

### Synopsis

Audit all HCP management clusters in the fleet to determine autoscaling migration readiness.

Outputs a flat table with autoscaling status across the entire fleet.

```
osdctl hcp autoscaling-audit [flags]
```

### Examples

```

  # Audit entire fleet
  osdctl hcp autoscaling-audit

  # Audit fleet with CSV output
  osdctl hcp autoscaling-audit --output csv > fleet-audit.csv

  # Show only clusters ready for migration
  osdctl hcp autoscaling-audit --show-only ready-for-migration

  # Show only clusters that need annotation removal
  osdctl hcp autoscaling-audit --show-only needs-removal

  # Show only clusters safe to remove override
  osdctl hcp autoscaling-audit --show-only safe-to-remove-override

  # Lower concurrency for slower networks
  osdctl hcp autoscaling-audit --concurrency 5
```

### Options

```
      --concurrency int    Number of management clusters to audit concurrently (default 10)
  -h, --help               help for autoscaling-audit
      --no-headers         Skip table headers in output
      --output string      Output format: text, json, yaml, csv (default "text")
      --show-only string   Filter output: needs-removal, ready-for-migration, safe-to-remove-override
```

### Options inherited from parent commands

```
      --as string                        Username to impersonate for the operation. User could be a regular user or a service account in a namespace.
      --cluster string                   The name of the kubeconfig cluster to use
      --context string                   The name of the kubeconfig context to use
      --insecure-skip-tls-verify         If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string                Path to the kubeconfig file to use for CLI requests.
      --request-timeout string           The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
  -s, --server string                    The address and port of the Kubernetes API server
      --skip-aws-proxy-check aws_proxy   Don't use the configured aws_proxy value
  -S, --skip-version-check               skip checking to see if this is the most recent release
```

### SEE ALSO

* [osdctl hcp](osdctl_hcp.md)	 - 

