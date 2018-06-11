# ![](assets/web/gitwebhookproxy-round-100px.png)  GitWebhookProxy

[![Go Report Card](https://goreportcard.com/badge/github.com/stakater/GitWebhookProxy?style=flat-square)](https://goreportcard.com/report/github.com/stakater/GitWebhookProxy)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/stakater/GitWebhookProxy)
[![Release](https://img.shields.io/github/release/stakater/GitWebhookProxy.svg?style=flat-square)](https://github.com/stakater/GitWebhookProxy/releases/latest)
[![GitHub tag](https://img.shields.io/github/tag/stakater/GitWebhookProxy.svg?style=flat-square)](https://github.com/stakater/GitWebhookProxy/releases/latest)
[![Docker Pulls](https://img.shields.io/docker/pulls/stakater/gitwebhookproxy.svg?style=flat-square)](https://hub.docker.com/r/stakater/GitWebhookProxy/)
[![Docker Stars](https://img.shields.io/docker/stars/stakater/gitwebhookproxy.svg?style=flat-square)](https://hub.docker.com/r/stakater/GitWebhookProxy/)
[![MicroBadger Size](https://img.shields.io/microbadger/image-size/jumanjiman/puppet.svg?style=flat-square)](https://microbadger.com/images/stakater/GitWebhookProxy)
[![MicroBadger Layers](https://img.shields.io/microbadger/layers/_/httpd.svg?style=flat-square)](https://microbadger.com/images/stakater/GitWebhookProxy)
[![license](https://img.shields.io/github/license/stakater/GitWebhookProxy.svg?style=flat-square)](LICENSE)

A proxy to let webhooks to reach a Jenkins instance running behind a firewall

## PROBLEM

Jenkins is awesome and matchless tool for both CI & CD; but unfortunately its a gold mine if left in wild with wide open access; so, we always want to put it behind a firewall. But when we put it behind firewall then webhooks don't work anymore and no one wants the pull based polling but rather prefer the build to start as soon as there is a commit!

## SOLUTION

This little proxy makes webhooks start working again!

### Supported Providers

Currently we support the following git providers out of the box:

* Github
* Gitlab

### Configuration

GitWebhookProxy can be configured by providing the following arguments either via command line or via environment variables:

| Parameter     | Description                                                                   | Default  | Example                                    |
|---------------|-------------------------------------------------------------------------------|----------|--------------------------------------------|
| listenAddress | Address on which the proxy listens.                                           | `:8080`  | `127.0.0.1:80`                             |
| upstreamURL   |          URL to which the proxy requests will be forwarded (required)         |          | `https://someci-instance-url.com/webhook/` |
| secret        | Secret of the Webhook API (required)                                          |          | `iamasecret`                               |
| provider      | Git Provider which generates the Webhook                                      | `github` | `github` or `gitlab`                       |
| allowedPaths  | Comma-Separated String List of allowed paths on the proxy                     |          | `/project` or `github-webhook/,project/`   |
| ignoredUsers  | Comma-Separated String List of users to ignore while proxying Webhook request |          | `someuser`                                 |

## DEPLOYING TO KUBERNETES

### Vanilla Manifests

You have to first clone or download the repository contents. The kubernetes deployment files are provided inside `deployments/kubernetes/manifests` folder.

Once you have the repo cloned, you can deploy GitwebhookProxy by running the following kubectl commands:

```bash
kubectl apply -f configmap.yaml -n <namespace>
kubectl apply -f secret.yaml -n <namespace>
kubectl apply -f deployment.yaml -n <namespace>
kubectl apply -f service.yaml -n <namespace>
```

*Note:* Make sure to update the `port` in deployment.yaml as well as service.yaml if you change the default `listenAddress` port.

### Helm Charts

Alternatively if you have configured helm on your cluster, you can deploy the controller via helm chart located under `deployment/kubernetes/chart` folder.

## Help

**Got a question?**
File a GitHub [issue](https://github.com/stakater/GitWebhookProxy/issues), or send us an [email](mailto:stakater@gmail.com).

### Talk to us on Slack
Join and talk to us on the #tools-gwp channel for discussing about GitWebhookProxy

[![Join Slack](https://stakater.github.io/README/stakater-join-slack-btn.png)](https://stakater-slack.herokuapp.com/)
[![Chat](https://stakater.github.io/README/stakater-chat-btn.png)](https://stakater.slack.com/messages/CAQ5A4HGD/)

## Contributing

### Bug Reports & Feature Requests

Please use the [issue tracker](https://github.com/stakater/GitWebhookProxy/issues) to report any bugs or file feature requests.

### Developing

PRs are welcome. In general, we follow the "fork-and-pull" Git workflow.

 1. **Fork** the repo on GitHub
 2. **Clone** the project to your own machine
 3. **Commit** changes to your own branch
 4. **Push** your work back up to your fork
 5. Submit a **Pull request** so that we can review your changes

NOTE: Be sure to merge the latest from "upstream" before making a pull request!

## Changelog

View our closed [Pull Requests](https://github.com/stakater/GitWebhookProxy/pulls?q=is%3Apr+is%3Aclosed).

## License

Apache2 © [Stakater](http://stakater.com)

## About

`GitWebhookProxy` is maintained by [Stakater][website]. Like it? Please let us know at <hello@stakater.com>

See [our other projects][community]
or contact us in case of professional services and queries on <hello@stakater.com>

  [website]: http://stakater.com/
  [community]: https://github.com/stakater/
