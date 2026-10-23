# kubectl user password manager glue

![CI status badge](https://github.com/chrisns/kubectl-passman/workflows/CI%20Pipeline/badge.svg)
![LICENSE](https://img.shields.io/github/license/chrisns/kubectl-passman)
![GitHub watchers](https://img.shields.io/github/watchers/chrisns/kubectl-passman?style)
![GitHub stars](https://img.shields.io/github/stars/chrisns/kubectl-passman)
![GitHub forks](https://img.shields.io/github/forks/chrisns/kubectl-passman)
![GitHub issues](https://img.shields.io/github/issues-raw/chrisns/kubectl-passman)
![GitHub closed issues](https://img.shields.io/github/issues-closed-raw/chrisns/kubectl-passman)
![GitHub pull requests](https://img.shields.io/github/issues-pr-raw/chrisns/kubectl-passman)
![GitHub closed pull requests](https://img.shields.io/github/issues-pr-closed-raw/chrisns/kubectl-passman)
![GitHub repo size](https://img.shields.io/github/repo-size/chrisns/kubectl-passman)
![GitHub contributors](https://img.shields.io/github/contributors/chrisns/kubectl-passman)
![GitHub last commit](https://img.shields.io/github/last-commit/chrisns/kubectl-passman)
[![Go Report Card](https://goreportcard.com/badge/github.com/chrisns/kubectl-passman)](https://goreportcard.com/report/github.com/chrisns/kubectl-passman)

 > :heavy_exclamation_mark: An easy way to store your kubernetes credentials in a keychain or password manager

### Does your `~/.kube/config` look like this:

```yaml
apiVersion: v1
kind: Config
users:
- name: my-prod-user
  user:
    token: <REAL TOKEN!>
- name: docker-desktop
  user:
    client-certificate-data: <REAL CERT!>
    client-key-data: <REAL PRIVATE KEY!>
```

## :scream: :scream: :scream: :scream:<br/><br/>Do you scold your parents :man_teacher:/:woman_teacher: for maintaining a `passwords.doc` on their desktop? <br/><br/> Then you need kubectl-passman!

## Works with (more coming)

Provider | Supports | Example command
--- | --- | ---
keychain | [Mac OS Keychain](https://support.apple.com/en-gb/guide/keychain-access/kyca1083/mac) <br> [GNOME Keyring](https://wiki.gnome.org/Projects/GnomeKeyring) <br> [Windows Credential Manager](http://blogs.msdn.com/b/visualstudioalm/archive/2015/12/08/announcing-the-git-credential-manager-for-windows-1-0.aspx) | `kubectl passman keychain [item] [token]`
1password | [1password](https://1password.com/) <br> requires [1password cli](https://1password.com/downloads/command-line/) | `kubectl passman 1password [item] [token]`
gopass | [gopass](https://www.gopass.pw/) | `kubectl passman gopass [item] [token]`
keepassxc | [KeePassXC](https://keepassxc.org/) <br> requires [keepassxc-cli](https://keepassxc.org/docs/KeePassXC_UserGuide#_command_line_tool) <br> on Linux you can also use `keychain`, see [KeePassXC](#keepassxc) | `kubectl passman keepassxc [item] [token]`

## Installation

```bash
# with krew (recommended)
kubectl krew install passman

# get a binary from https://github.com/chrisns/kubectl-passman/releases/latest
# place it in PATH and make sure it's called kubectl-passman

# use go to get the most recent
go install github.com/chrisns/kubectl-passman
```

## Usage

You need to JSON encode the credentials so that should look something like:

```json
{"token":"00000000-0000-0000-0000-000000000000"}
```

or just the token on its own, for example an OIDC `id-token` from Dex:

```
eyJhbGciOiJSUzI1NiIs...
```

or the `id-token` key from the `auth-provider` config of your kube config:

```json
{"id-token":"eyJhbGciOiJSUzI1NiIs..."}
```

A bare token can only have letters, digits and `-._~+/=`. Passman removes white space around it.
Passman ignores `refresh-token`, `client-id` and `client-secret`.
Passman does not refresh an `id-token`. When the token expires, store a new one.
For automatic refresh, use [kubelogin](https://github.com/int128/kubelogin).

or for a key pair:

```json
{
  "clientCertificateData":"-----BEGIN REAL CERTIFICATE-----\nMIIC9DCCA.......-----END CERTIFICATE-----",
  "clientKeyData":"-----BEGIN REAL RSA PRIVATE KEY-----\nMIIE......-----END REAL RSA PRIVATE KEY-----"
}
```

or for a key pair from your kube config:

```json
{
  "client-certificate-data":"LS0tLS1CRU...LS0tCg==",
  "client-key-data":"LS0tLS1CRU...LS0tLS0K"
}
```

If they are already in your kube config, you could retrieve them with something like:

```bash
kubectl config view --raw -o json | jq '.users[] | select(.name=="kubectl-prod-user") | .user' -c
```

### Write it to the password manager

```bash
kubectl passman keychain kubectl-prod-user '[token]'
# or
kubectl passman 1password kubectl-prod-user '[token]'

## so should look like:
kubectl passman 1password kubectl-prod-user '{"token":"00000000-0000-0000-0000-000000000000"}'
# or
kubectl passman 1password kubectl-prod-user '{"client-certificate-data":"...BASE64_ENCODE...","client-key-data":"...BASE64_ENCODE..."}'
```

Then add it to the `~/.kube/config`:

```bash
kubectl config set-credentials \
  kubectl-prod-user \
 --exec-api-version=client.authentication.k8s.io/v1beta1 \
 --exec-command=kubectl-passman \
 --exec-arg=keychain \ # or 1password
 --exec-arg=kubectl-prod-user # name of [item-name] you used when you wrote to the password manager
```

### KeePassXC

A KeePass database is a file, so nothing keeps it unlocked between runs.
There are two ways to avoid a password prompt on every `kubectl` command.

**Linux: use the `keychain` provider with the Secret Service.**
KeePassXC can act as the [Secret Service](https://keepassxc.org/docs/KeePassXC_UserGuide#_secret_service_integration) provider.
Enable it in KeePassXC under Tools > Settings > Secret Service Integration, and expose a database group.
Then use `kubectl passman keychain [item] [token]`.
You unlock the database once in KeePassXC, and `kubectl-passman` reads from the unlocked database.

**Any platform: use the `keepassxc` provider.**
It runs `keepassxc-cli` and stores the credentials in the `Password` field of the entry.
Set these environment variables:

Variable | Meaning
--- | ---
`KUBECTL_PASSMAN_KEEPASSXC_DB` | Path to the `.kdbx` file (required)
`KUBECTL_PASSMAN_KEEPASSXC_KEYFILE` | Path to a key file (optional)
`KUBECTL_PASSMAN_KEEPASSXC_PASSWORD` | Database password (optional, see the warning below)
`KUBECTL_PASSMAN_KEEPASSXC_NO_PASSWORD` | Set to any non-empty value for a database that has a key file and no password

If you set no password, `kubectl-passman` prompts for it on the terminal each time.
The prompt needs `/dev/tty`, so on Windows you must set `KUBECTL_PASSMAN_KEEPASSXC_PASSWORD` or use a key file with `KUBECTL_PASSMAN_KEEPASSXC_NO_PASSWORD`.
`kubectl` does not pass your shell variables to the plugin, so set them with `--exec-env` as shown below.
Do not put the database password in a file that other users can read.
You can also write `keepass` instead of `keepassxc`.

```bash
export KUBECTL_PASSMAN_KEEPASSXC_DB=~/Passwords.kdbx # needed to write the item
kubectl passman keepassxc kubectl-prod-user '{"token":"00000000-0000-0000-0000-000000000000"}'
kubectl config set-credentials \
  kubectl-prod-user \
 --exec-api-version=client.authentication.k8s.io/v1beta1 \
 --exec-command=kubectl-passman \
 --exec-arg=keepassxc \
 --exec-arg=kubectl-prod-user \
 --exec-env=KUBECTL_PASSMAN_KEEPASSXC_DB=$HOME/Passwords.kdbx
```

## Build

``` bash
go build
```
> Note: kubectl-passman will build slightly differently on Darwin (Mac OS) to other operation systems because it uses the [go-keychain](https://github.com/keybase/go-keychain) library that needs libraries that only exist on a mac so that it can natively talk to the keychain. When compiling for other operating systems you'll get [go-keyring](https://github.com/zalando/go-keyring) instead but I've abstracted to make the interactions the same.

## Contributing

I :heart: contributions, it'd be great if you could add support for your favourite password manager, work on something from the [TODO](#TODO) or any open issues as a priority, but anything else that takes your fancy too is great, though best to raise an issue to discuss before investing time into it.

## Thanks

<a href="https://www.atlascloud.ai/"><img src="https://www.atlascloud.ai/logo.svg" alt="Atlas Cloud" height="20"></a> sponsors the cloud infrastructure for this project.
