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

or for a key pair:

```json
{
  "clientCertificateData":"-----BEGIN REAL CERTIFICATE-----\nMIIC9DCCA.......-----END CERTIFICATE-----",
  "clientKeyData":"-----BEGIN REAL RSA PRIVATE KEY-----\nMIIE......-----END REAL RSA PRIVATE KEY-----"
}
```

If they are already in your kube config, you could retrieve them with something like:

```bash
kubectl config view -o json | jq '.users[] | select(.name=="kubectl-prod-user") | .user' -c
```

### Write it to the password manager

```bash
kubectl passman keychain kubectl-prod-user '[token]'
# or
kubectl passman 1password kubectl-prod-user '[token]'

## so should look like:
kubectl passman 1password kubectl-prod-user '{"token":"00000000-0000-0000-0000-000000000000"}'

```

Then add it to the `~/.kube/config`:

```bash
kubectl config set-credentials \
  kubectl-prod-user \
 --exec-api-version=client.authentication.k8s.io/v1beta \
 --exec-command=passman \
 --exec-arg=keychain \ # or 1password
 --exec-arg=kubectl-prod-user # name of [item-name] you used when you wrote to the password manager
```

## Build

``` bash
go build
```
> Note: kubectl-passman will build slightly differently on Darwin (Mac OS) to other operation systems because it uses the [go-keychain](https://github.com/keybase/go-keychain) library that needs libraries that only exist on a mac so that it can natively talk to the keychain. When compiling for other operating systems you'll get [go-keyring](https://github.com/zalando/go-keyring) instead but I've abstracted to make the interactions the same.

## Contributing

I :heart: contributions, it'd be great if you could add support for your favourite password manager, work on something from the [TODO](#TODO) or any open issues as a priority, but anything else that takes your fancy too is great, though best to raise an issue to discuss before investing time into it.
