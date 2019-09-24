# kubectl user password manager glue

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

## :scream: :scream: :scream: :scream:<br/><br/>Do you scold your parents :man_teacher:/:woman_teacher: for maintaining a `passwords.doc` on their desktop?

## Then you need kubectl-passman!

> a suggestion of a better name is very welcome, it's not too late to change!

What this `kubectl` plugin does is glue your kubectl config to a common password manager.

## Installation

### Mac OS X/Unix/Linux

```bash
go install github.com/chrisns/kubectl-passman
```

### Windows

```powershell
TODO:
```

## Usage

### macOS Keychain

You need to JSON encode the credentials so that should look something like:

```json
{"token":"some-token"}
```

or for a key pair:

```json
{
  "clientCertificateData":"-----BEGIN CERTIFICATE-----\nMIIC9DCCA.......-----END CERTIFICATE-----",
  "clientKeyData":"-----BEGIN RSA PRIVATE KEY-----\nMIIE......-----END RSA PRIVATE KEY-----"
}
```

You then place this in a keychain item, call it whatever you like but keep the account name and item name the same.

![Screenshot of adding a keypair](resources/osxkeychain-keypair.png)
![Screenshot of adding a token](resources/osxkeychain-token.png)

Then add it to the `~/.kube/config`:

```yaml
apiVersion: v1
kind: Config
users:
- name: my-prod-user
    user:
      exec:
        command: "kubectl-passman"
        apiVersion: "client.authentication.k8s.io/v1beta1"
        args:
          - keychain
          - kubectl-prod-user
```

## Compiling

``` bash
go build
```

## Contributing

I :heart: contributions, it'd be great if you could add support for your favourite password manager, work on something from the [TODO](#TODO) or any open issues as a priority, but anything else that takes your fancy too is great, though best to raise an issue to discuss before investing time into it.

## TODO

- [x] rename project ~~k8s-user-passmanager~~ kubectl-passman to math a compatible named binary
- [x] skeleton readme doc
- [ ] retrieve tokens
  - [x] from osx keychain
  - [ ] from 1Password
    - [ ] in os x
    - [ ] in windows
    - [ ] in *nix
- [ ] retrieve cert key/pair
  - [x] from osx keychain
  - [ ] from 1Password
    - [ ] in os x
    - [ ] in windows
    - [ ] in *nix
- [ ] validate that kubectl is official signed build before executing
- [ ] CI/CD/CD
  - [ ] test
  - [ ] build binaries
    - [ ] osx
    - [ ] *nix
    - [ ] windows
    - [ ] x86?
    - [ ] arm?
    - [ ] docker?
  - [ ] publish binaries to github releases
  - [ ] cli interface for abstracting creating new credentials in your password manager e.g. `kubectl passman keychain create [item name] --token=[my token]`
