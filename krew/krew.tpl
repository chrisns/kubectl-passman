apiVersion: krew.googlecontainertools.github.com/v1alpha2
kind: Plugin
metadata:
  name: passman
spec:
  version: "{{ env "VERSION"}}"
  platforms:
    - selector:
        matchLabels:
          os: darwin
          arch: amd64
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "REF"}}/kubectl-passman-darwin-amd64.zip
      sha256: "{{.kubectl_passman_darwin_amd64}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-darwin-amd64
          to: kubectl-passman

    - selector:
        matchLabels:
          os: darwin
          arch: 386
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "REF"}}/kubectl-passman-darwin-386.zip
      sha256: "{{.kubectl_passman_darwin_386}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-darwin-386
          to: kubectl-passman

    - selector:
        matchLabels:
          os: freebsd
          arch: 386
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "REF"}}/kubectl-passman-freebsd-386.zip
      sha256: "{{.kubectl_passman_freebsd_386}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-freebsd-386
          to: kubectl-passman

    - selector:
        matchLabels:
          os: freebsd
          arch: amd64
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "REF"}}/kubectl-passman-freebsd-amd64.zip
      sha256: "{{.kubectl_passman_freebsd_386}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-freebsd-amd64
          to: kubectl-passman

    - selector:
        matchLabels:
          os: freebsd
          arch: arm
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "REF"}}/kubectl-passman-freebsd-arm.zip
      sha256: "{{.kubectl_passman_freebsd_arm}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-freebsd-arm
          to: kubectl-passman

    - selector:
        matchLabels:
          os: linux
          arch: arm
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "REF"}}/kubectl-passman-linux-arm.zip
      sha256: "{{.kubectl_passman_linux_arm}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-linux-arm
          to: kubectl-passman

    - selector:
        matchLabels:
          os: linux
          arch: 386
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "REF"}}/kubectl-passman-linux-386.zip
      sha256: "{{.kubectl_passman_linux_386}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-linux-386
          to: kubectl-passman

    - selector:
        matchLabels:
          os: linux
          arch: amd64
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "REF"}}/kubectl-passman-linux-amd64.zip
      sha256: "{{.kubectl_passman_linux_amd64}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-linux-amd64
          to: kubectl-passman

    - selector:
        matchLabels:
          os: linux
          arch: arm64
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "REF"}}/kubectl-passman-linux-arm64.zip
      sha256: "{{.kubectl_passman_linux_arm64}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-linux-arm64
          to: kubectl-passman

    - selector:
        matchLabels:
          os: netbsd
          arch: arm64
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "REF"}}/kubectl-passman-netbsd-arm64.zip
      sha256: "{{.kubectl_passman_netbsd_arm64}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-netbsd-arm64
          to: kubectl-passman

    - selector:
        matchLabels:
          os: netbsd
          arch: arm
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "REF"}}/kubectl-passman-netbsd-arm.zip
      sha256: "{{.kubectl_passman_netbsd_arm}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-netbsd-arm
          to: kubectl-passman

    - selector:
        matchLabels:
          os: netbsd
          arch: 386
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "REF"}}/kubectl-passman-netbsd-386.zip
      sha256: "{{.kubectl_passman_netbsd_386}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-netbsd-386
          to: kubectl-passman

    - selector:
        matchLabels:
          os: netbsd
          arch: amd64
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "REF"}}/kubectl-passman-netbsd-amd64.zip
      sha256: "{{.kubectl_passman_netbsd_amd64}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-netbsd-amd64
          to: kubectl-passman

    - selector:
        matchLabels:
          os: openbsd
          arch: amd64
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "REF"}}/kubectl-passman-openbsd-amd64.zip
      sha256: "{{.kubectl_passman_openbsd_amd64}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-openbsd-amd64
          to: kubectl-passman

    - selector:
        matchLabels:
          os: openbsd
          arch: arm
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "REF"}}/kubectl-passman-openbsd-arm.zip
      sha256: "{{.kubectl_passman_openbsd_arm}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-openbsd-arm
          to: kubectl-passman

    - selector:
        matchLabels:
          os: openbsd
          arch: arm64
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "REF"}}/kubectl-passman-openbsd-arm64.zip
      sha256: "{{.kubectl_passman_openbsd_arm64}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-openbsd-arm64
          to: kubectl-passman

    - selector:
        matchLabels:
          os: openbsd
          arch: 386
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "REF"}}/kubectl-passman-openbsd-386.zip
      sha256: "{{.kubectl_passman_openbsd_386}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-openbsd-386
          to: kubectl-passman

    - selector:
        matchLabels:
          os: plan9
          arch: 386
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "REF"}}/kubectl-passman-plan9-386.zip
      sha256: "{{.kubectl_passman_plan9_386}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-plan9-386
          to: kubectl-passman

    - selector:
        matchLabels:
          os: plan9
          arch: amd64
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "REF"}}/kubectl-passman-plan9-amd64.zip
      sha256: "{{.kubectl_passman_plan9_amd64}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-plan9-amd64
          to: kubectl-passman

    - selector:
        matchLabels:
          os: plan9
          arch: arm
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "REF"}}/kubectl-passman-plan9-arm.zip
      sha256: "{{.kubectl_passman_plan9_arm}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-plan9-arm
          to: kubectl-passman

    - selector:
        matchLabels:
          os: solaris
          arch: amd64
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "REF"}}/kubectl-passman-solaris-amd64.zip
      sha256: "{{.kubectl_passman_solaris_amd64}}"
      bin: "./kubectl-passman"
      files:
        - from: kubectl-passman-solaris-amd64
          to: kubectl-passman

    - selector:
        matchLabels:
          os: windows
          arch: amd64
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "REF"}}/kubectl-passman-windows-amd64.zip
      sha256: "{{.kubectl_passman_windows_amd64}}"
      bin: "./kubectl-passman.exe"
      files:
        - from: kubectl-passman-windows-amd64.exe
          to: kubectl-passman.exe

    - selector:
        matchLabels:
          os: windows
          arch: 386
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "REF"}}/kubectl-passman-windows-386.zip
      sha256: "{{.kubectl_passman_windows_386}}"
      bin: "./kubectl-passman.exe"
      files:
        - from: kubectl-passman-windows-386.exe
          to: kubectl-passman.exe

    - selector:
        matchLabels:
          os: windows
          arch: arm
      uri: https://github.com/chrisns/kubectl-passman/releases/download/{{env "REF"}}/kubectl-passman-windows-arm.zip
      sha256: "{{.kubectl_passman_windows_arm}}"
      bin: "./kubectl-passman.exe"
      files:
        - from: kubectl-passman-windows-arm.exe
          to: kubectl-passman.exe

  shortDescription: kubectl plugin that aspires to provide the missing link/glue between common password managers and kubectl
  homepage: https://github.com/chrisns/kubectl-passman
  caveats: |
    This plugin needs a usable keychain or password manager, see the docs https://github.com/chrisns/kubectl-passman
  description: |
    An easy way to store your kubernetes credentials in 1password or Mac OS Keychain (more password managers coming soon)
