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
# :scream: :scream: :scream: :scream:

### Do you scold your parents :man_teacher:/:woman_teacher: for maintaining a `passwords.doc` on their desktop?

## Then you need kubectl-passman :thumbsup:!

> a suggestion of a better name is very welcome, it's not too late to change!

What this `kubectl` plugin does is glue your kubectl config to a common password manager.

## Installation

### Mac OS X

```bash
TODO:
```
### Windows

```powershell
TODO:
```

### Unix

```bash
TODO:
```

### Usage

#### Adding your keys/tokens to your password manager

#### What your `~/.kube/config` should look like

```yaml
TODO:
```

## Compiling
``` bash
TODO: 
```

## Contributing
I :heart: contributions, it'd be great if you could add support for your favourite password manager, work on something from the [TODO](#TODO) or any open issues as a priority, but anything else that takes your fancy too is great, though best to raise an issue to discuss before investing time into it.

## TODO

- [x] rename project ~~k8s-user-passmanager~~ kubectl-passman to math a compatible named binary
- [x] skeleton readme doc
- [ ] store and retrieve tokens
  - [ ] from osx keychain
  - [ ] from 1Password
    - [ ] in os x
    - [ ] in windows
    - [ ] in *nix
- [ ] store and retrieve private keys
  - [ ] from osx keychain
  - [ ] from 1Password
    - [ ] in os x
    - [ ] in windows
    - [ ] in *nix
- [ ] validate that kubectl is official signed build before executing
- [ ] CI/CD/CD
  - [ ] lint
  - [ ] test
  - [ ] build
    - [ ] osx
    - [ ] *nix
    - [ ] windows
    - [ ] docker?
    - [ ] arm?
  - [ ] publish to github releases
