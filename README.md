# `hdfs-to-local`

This is an experimental project.

Performs a recursive directory copy from HDFS host to local storage. Does not
require any Hadoop binaries since it uses Go native implemented HDFS client
library.

## Requirements

* [`go`](https://golang.org/dl/) + [`glide`](https://glide.sh/) **OR**
* [`docker`](https://www.docker.com/get-docker) +
  [`docker-compose`](https://docs.docker.com/compose/install/) that is able to
  run Compose v2 file

## How to Build

### Native Build

Install the go dependencies for this project into `vendor/`:

```bash
glide install
```

Build and you are done:

```bash
go build
```

The compiled executable will be located at the repository root directory, and is
named `hdfs-to-local`.

For fully statically linked executable, use the following command instead:

```bash
CGO_ENABLED=0 go build
```

### Docker Build

This alternative method is recommended if you do not wish to up `go`, `glide` or
the environment variable `GOPATH`, or prefer not to have this repository files
not within the `GOPATH`.

Run the following for full compilation:

```bash
docker-compose run -u $UID:`id -g` all
```

Note that `` -u $UID:`id -g` `` is optional, but useful to perform the
compilation under your current user's UID and GID.

The compiled executable will be located at the repository root directory, and is
named `hdfs-to-local`. The executable is always fully statically linked.

The following commands are available to run for `docker-compose`:

* `all`
  * Performs `glide install`, followed by `go build`.
  * e.g. `` docker-compose run -u $UID:`id -g` all ``
* `install`
  * Performs only `glide install`.
  * e.g. `` docker-compose run -u $UID:`id -g` install ``
* `build`
  * Performs only `go build`.
  * e.g. `` docker-compose run -u $UID:`id -g` build ``
* `clean`
  * Performs `go clean`.
  * e.g. `` docker-compose run -u $UID:`id -g` clean ``
* `update`
  * Performs `glide update`. Not recommended since it updates the dependencies.
  * e.g. `` docker-compose run -u $UID:`id -g` update ``

## How to Run

Assuming HDFS server is running on port 9000 locally, and there is a directory
at `/data`:

```bash
./hdfs-to-local --conf config/example.toml
```

For more program argument details:

```bash
./hdfs-to-local --help
```
