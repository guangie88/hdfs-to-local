# `hdfs-to-local`

This is still an experimental project.

Performs a recursive directory copy from HDFS host to local storage. Does not
require any Hadoop binaries since it uses Go native implemented HDFS client
library.

## Requirements

* [`glide`](https://glide.sh/)

## How to Build

Install the go dependencies for this project into `vendor/`:

```bash
glide install
```

Build and you are done:

```bash
go build
```

The compiled executable is named `hdfs-to-local`.

For pure statically linked executable, use the following command instead:

```bash
CGO_ENABLED=0 go build
```

## How to Run

Assuming HDFS server is running on port 8020 locally, and there is a directory
at `/data`:

```bash
./hdfs-to-local --host localhost:8020 --root /data
```

For more program argument details:

```bash
./hdfs-to-local --help
```
