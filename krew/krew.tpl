apiVersion: krew.googlecontainertools.github.com/v1alpha2
kind: Plugin
metadata:
  name: passman
spec:
  version: "{{ if not (hasPrefix "v" (env "VERSION")) }}v0.0.0-{{ end }}{{ env "VERSION" }}"
  platforms:












