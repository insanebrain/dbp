# DBP (DOCKER BUILD PARENT)
DBP build multiple images with dependencies with other images. A multi-build is detect by change of a commit.
DBP generate also a document with data specified in each directory.

## Usage
```help
dbp
usage: dbp [<flags>] <command> [<args> ...]

A command-line dbp helper.

Flags:
      --help           Show context-sensitive help (also try --help-long and --help-man).
  -c, --config=CONFIG  Define specific configuration file
  -p, --path="."       Define execution path

Commands:
  help [<command>...]
    Show help.

  build dirty
    Build docker images for dirty repo

  build commit <commit>
    Build docker images for specific commit

  list
    List all images of directory

  generate dirty
    Generate readme of images for dirty repo

  generate commit <commit>
    Generate readme of images for specific commit

  generate all
    Generate readme of images for all

  generate index
    Generate a readme index

  version
    Display version
```
## Requirements

* Docker
* GIT

## Setup

### Install go-bindata

```
go get -u github.com/go-bindata/go-bindata
```

### Generate templates

```
go generate
```

## Config file dbp.yml
See documentation example for image configuration [here](doc/examples/dbp.yml).

See documentation example for dbp configuration [here](doc/examples/config.yml).

## Bash/ZSH Shell Completion

### Bash

Add in your ~/.bashrc :
```
eval "$(dbp --completion-script-bash)"
```

### ZSH

Add in your ~/.zshrc :
```
eval "$(dbp --completion-script-zsh)"
```

## Docker

docker run -it -v /var/run/docker.sock:/var/run/docker.sock -v ${PWD}:/build insanebrain/dbp:master dbp build --push

## Compatibility
TODO

## Contributing

Refer to [CONTRIBUTING.md](CONTRIBUTING.md)

## License

Apache License 2.0, see [LICENSE](LICENSE.md).
