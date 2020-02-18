{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "gitwebhookproxy.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" | lower -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "gitwebhookproxy.fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "gitwebhookproxy.labels.selector" -}}
{{- if .Values.gitWebhookProxy.useCustomName -}}
app: {{ .Values.gitWebhookProxy.customName }}
{{- else -}}
app: {{ template "gitwebhookproxy.name" . }}
{{- end }}
group: {{ .Values.gitWebhookProxy.labels.group }}
provider: {{ .Values.gitWebhookProxy.labels.provider }}
{{- end -}}

{{- define "gitwebhookproxy.labels.stakater" -}}
{{ template "gitwebhookproxy.labels.selector" . }}
version: {{ .Values.gitWebhookProxy.labels.version }}
{{- end -}}

{{- define "gitwebhookproxy.labels.chart" -}}
chart: "{{ .Chart.Name }}-{{ .Chart.Version }}"
release: {{ .Release.Name | quote }}
heritage: {{ .Release.Service | quote }}
{{- end -}}

{{/*
Return the appropriate apiVersion for deployment.
*/}}
{{- define "deployment.apiVersion" -}}
{{- if semverCompare ">=1.9-0" .Capabilities.KubeVersion.GitVersion -}}
{{- print "apps/v1" -}}
{{- else -}}
{{- print "extensions/v1beta1" -}}
{{- end -}}
{{- end -}}