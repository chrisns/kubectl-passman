apiVersion: krew.googlecontainertools.github.com/v1alpha2
kind: Plugin
metadata:
  name: passman
spec:
  version: "{{ if not (hasPrefix "v" (env "VERSION")) }}v0.0.0-{{ end }}{{ env "VERSION" }}"
  platforms:
    - selector:
        matchLabels:
          os: darwin
          arch: amd64
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "VERSION"}}/kubectl-passman-darwin-amd64.zip
      sha256: "{{.kubectl_passman_darwin_amd64}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-darwin-amd64
          to: kubectl-passman

    - selector:
        matchLabels:
          os: darwin
          arch: 386
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "VERSION"}}/kubectl-passman-darwin-386.zip
      sha256: "{{.kubectl_passman_darwin_386}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-darwin-386
          to: kubectl-passman

    - selector:
        matchLabels:
          os: linux
          arch: arm
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "VERSION"}}/kubectl-passman-linux-arm.zip
      sha256: "{{.kubectl_passman_linux_arm}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-linux-arm
          to: kubectl-passman

    - selector:
        matchLabels:
          os: linux
          arch: 386
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "VERSION"}}/kubectl-passman-linux-386.zip
      sha256: "{{.kubectl_passman_linux_386}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-linux-386
          to: kubectl-passman

    - selector:
        matchLabels:
          os: linux
          arch: amd64
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "VERSION"}}/kubectl-passman-linux-amd64.zip
      sha256: "{{.kubectl_passman_linux_amd64}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-linux-amd64
          to: kubectl-passman

    - selector:
        matchLabels:
          os: windows
          arch: amd64
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "VERSION"}}/kubectl-passman-windows-amd64.zip
      sha256: "{{.kubectl_passman_windows_amd64}}"
      bin: "./kubectl-passman.exe"
      files:
        - from: kubectl-passman-windows-amd64.exe
          to: kubectl-passman.exe

    - selector:
        matchLabels:
          os: windows
          arch: 386
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "VERSION"}}/kubectl-passman-windows-386.zip
      sha256: "{{.kubectl_passman_windows_386}}"
      bin: "./kubectl-passman.exe"
      files:
        - from: kubectl-passman-windows-386.exe
          to: kubectl-passman.exe

  shortDescription: Store kubeconfig credentials in keychains or password managers
  homepage: https://github.com/chrisns/kubectl-passman
  caveats: |
    This plugin needs a usable keychain or password manager
    See usage docs https://github.com/chrisns/kubectl-passman
  description: |
    An easy way to store your kubernetes credentials in 1password,
    GNOME Keyring or Mac OS Keychain
    (more password managers coming soon)
