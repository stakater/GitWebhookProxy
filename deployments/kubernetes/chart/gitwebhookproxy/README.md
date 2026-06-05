git diff# GitWebhookProxy Helm Chart

This Helm chart deploys GitWebhookProxy on Kubernetes with support for GCP-specific features.

## GCP Features

The chart now supports the following GCP-specific features:

### FrontendConfig

Enable and configure a GCP FrontendConfig resource:

```yaml
gitWebhookProxy:
  ingress:
    gcp:
      enabled: true
      name: "my-frontend-config"  # Optional custom name
      sslPolicy:
        enabled: true
        name: "modern-ssl-policy"
      redirects:
        enabled: true
        type: PERMANENT_REDIRECT
        responseCode: MOVED_PERMANENTLY_DEFAULT
```

### BackendConfig

Configure a named backend with health checks:

```yaml
gitWebhookProxy:
  ingress:
    gcp:
      backend:
        enabled: true
        name: "my-backend-config"  # Optional custom name
        healthCheck:
          enabled: true
          checkIntervalSec: 30
          timeoutSec: 5
          healthyThreshold: 2
          unhealthyThreshold: 3
          type: HTTP
          port: 8080
          requestPath: "/health"