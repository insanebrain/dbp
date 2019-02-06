# DBP (Docker Build Parent)

DBP build multiple images with dependencies with other images.

A multi-build is detect by change of a commit.

DBP generate also a document with data specified in each directory.

[![Build Status](https://travis-ci.com/insanebrain/dbp.svg?branch=master)](https://travis-ci.com/insanebrain/dbp)

## Usage

```help
dbp --help
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

* Docker (18.06+)
* Git (2.17+)

## Development

* Install go-bindata:

  ```bash
  go get -u github.com/go-bindata/go-bindata/...
  ```

* Generate assets:

  ```bash
  go generate
  ```

* Build the project:

  ```bash
  go build .
  ```

* Run tests:

  ```bash
  go test ./...
  ```

## Production

* Run dbp with docker:

  ```bash
  docker run -it -v /var/run/docker.sock:/var/run/docker.sock -v ${PWD}:/build insanebrain/dbp:${VERSION} dbp --help
  ```

* Install dbp with deb (may require sudo privileges):

  ```bash
  curl -L "https://github.com/insanebrain/dbp/releases/download/${VERSION}/dbp_$(uname -m).deb" -o dbp_$(uname -m).deb
  dpkg -i dbp_$(uname -m).deb
  rm dbp_$(uname -m).deb
  ```

* Install binary to custom location:

  ```bash
  curl -L "https://github.com/insanebrain/dbp/releases/download/${VERSION}/dbp_$(uname -s)_$(uname -m)" -o ${DESTINATION}/dbp
  chmod +x ${DESTINATION}/dbp
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

## Contributing

Refer to [CONTRIBUTING.md](CONTRIBUTING.md)

## License

MIT License, see [LICENSE](LICENSE.md).
