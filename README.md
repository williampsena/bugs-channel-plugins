# BugsChannel (Plugins)

![bugs channel logo](./images/logo.png)

![workflow](https://github.com/williampsena/bugs-channel-plugins/actions/workflows/main.yml/badge.svg)

This repository provides plugin integrations such as Sentry, Honeybadger, Rollbar, and others, which add value to the project.

> I started the project with Elixir, but I'm switching to Go to keep things as simple and productive as possible 😅.

I decided to begin this project with the goal of making error handling as simple as possible.
I use [Sentry](https://sentry.io) and [Honeybadger](https://www.honeybadger.io), and both tools are fantastic for quickly tracking down issues. However, the purpose of this project is not to replace them, but rather to provide a simple solution for you to run on premise that is easy and has significant features.

# Requirements

```shell
go install golang.org/x/pkgsite/cmd/pkgsite@latest
go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
```

# Challenges
## Done 👌

- Handle Sentry events from their SDKs
- Check for the presence of authentication keys
- In db-less mode, define yaml as an option
- Identify the project by the requested authentication keys

## TODO

- Implement the rate-limit strategy
- Generate documentation with pkgsite
- Handle Honeybadger events from their SDKs
- Handle Rollbar events from their SDKs

# Running project

The command below starts a web application that listens on port 4001 by default.

```shell
# verbose mode
cp .env.sample .env
make dev
```

The project listens on port 4001 (sentry). At the moment, just Sentry had been set up and you could test the following steps.

- Create a config file named `config.yml` to run as **dbless** mode.

```shell
cp fixtures/settings/config.yml .config/config.yml
```

- Create a file named `main.py`.

```python
import sentry_sdk

sentry_sdk.init(
    "http://key@localhost:4001/1",
    traces_sample_rate=1.0,
)

raise ValueError("Error SDK")
```

- Install python packages

```shell
# using venv
python -m venv .env
. .env/bin/activate
pip install sentry-sdk

# without venv
pip install --user sentry-sdk
```

- Now you can run project

```shell
python main.py
```

# Tests

```shell
make test
```

# Vulnerabilities

```shell
make vulns-check
```

# Docs

```shell
make docs
```
