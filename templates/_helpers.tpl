{{- define "schemahero.fullname" -}}
{{ .Release.Name }}
{{- end -}}

{{- define "schemahero.labels" -}}
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version }}
{{ include "schemahero.selectorLabels" . }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{- define "schemahero.selectorLabels" -}}
app.kubernetes.io/name: {{ include "schemahero.fullname" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
control-plane: schemahero
{{- end -}}

{{- define "schemahero.webhookSecret" -}}
{{ include "schemahero.fullname" . }}-webhook
{{- end -}}
