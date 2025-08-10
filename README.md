# bqdo

bqdo is a simple CLI for running BigQuery SQL pipelines from a directory of `.sql` files. It reads a `bqdo.toml` configuration, applies Go `text/template` variables, and executes each SQL file in alphabetical order. It is designed for lightweight, repeatable data workflows.

## Features
- Read `bqdo.toml` from the current directory (or `-c` path)
- Override config values with CLI flags
- Optional service account impersonation
- Recursively discover and run `.sql` files in alphabetical order
- Go `text/template` support using values from the `[vars]` section
- `--dry-run` to validate without executing

## Requirements
- Go `1.24+`
- Optional: `mise` for project tasks (recommended for building)
- Auth to Google Cloud (ADC via `gcloud auth application-default login`, or `GOOGLE_APPLICATION_CREDENTIALS` pointing to a service account JSON)

## Download & Build
Clone the repo and build using `mise` (recommended):

```sh
git clone https://github.com/wintermi/bqdo.git
cd bqdo
mise run build
```

This produces the binary at `bin/bqdo`.

Alternatively, build directly with Go:

```sh
CGO_ENABLED=0 go build -o bin/bqdo ./cmd/bqdo
```

## Quick Start
1) Initialize a new pipeline in your project directory:

```sh
bin/bqdo init
```

This creates a `bqdo.toml`. Edit it to set `directory`, `dataset`, and other values.

2) Add your SQL files under the configured `directory` (e.g., `sql/`). Files will execute in alphabetical order.

3) Run the pipeline:

```sh
bin/bqdo run -c bqdo.toml -p YOUR_PROJECT_ID -d YOUR_DATASET -l US
```

Use `--dry-run` to validate queries without executing:

```sh
bin/bqdo run -c bqdo.toml -p YOUR_PROJECT_ID -d YOUR_DATASET -l US --dry-run
```

To impersonate a service account:

```sh
bin/bqdo run -c bqdo.toml -p YOUR_PROJECT_ID -d YOUR_DATASET -l US \
  --impersonate-service-account your-sa@your-project.iam.gserviceaccount.com
```

## Examples
A minimal runnable example is provided in `examples/simple`. See its README for details.

## License
Apache 2.0. See `LICENSE`.
