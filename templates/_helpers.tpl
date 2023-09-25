{{- define "schemahero.fullname" -}}
{{ .Release.Name }}
{{- end -}}

{{- define "schemahero.labels" -}}
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version }}
{{ include "schemahero.selectorLabels" . }}
{{/* If image is passed with @sha256 part, we don't want it in labels */}}
app.kubernetes.io/version: {{ index (regexSplit "@sha256:" (.Values.image.tag | default .Chart.AppVersion) -1) 0 | quote }}
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
