{{- define "istio-ratelimit-operator.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "istioratelimitoperator.labels" -}}
{{- with .Values.extraLabels }}
{{ toYaml . }}
{{- end }}
{{- end }}
