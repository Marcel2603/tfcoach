# Docs

This folder contains the documenation of tfcoach

## Run locally

```shell
python3 -m venv .venv
source .venv/bin/activate
pip3 install -r requirements.txt
```

After that you can interact with mkdocs.

* `mkdocs serve` - Start the live-reloading docs server.
* `mkdocs build` - Build the documentation site.
* `mkdocs -h` - Print help message and exit.

### Run inside docker

```shell
docker build . -t mkdocs:tfcoach
docker run -p 8000:8000 mkdocs:tfcoach
```