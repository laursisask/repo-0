# Contrast Assess Count License Usage

Script to count licensed Contrast Assess Applications across environments, de-duplicating them by name, language, and metadata.
Intended for use with Prometheus.

A total unique application count metric is emitted, as well as used license counts for each environment.

## Requirements
- Python 3.10 (other versions _may_ work but are untested)
- Ability to install Python libraries from `requirements.txt`

## Setup
You can run this script locally with a Python install, or, in a container with the provided `Dockerfile`

### Container use

#### Local build
```bash
docker build . --tag contrast-count-assess-licenses # Build the container
docker run -it -v $PWD/config.json:/usr/src/app/config.json contrast-count-assess-licenses <...args...> # Run the container
```

### Local use
Use of a virtual environment is encouraged
```bash
python3 -m venv venv # Create the virtual environment
. venv/bin/activate # Activate the virtual environment
pip3 install -r requirements.txt # Install dependencies
python3 contrast_application_licenses.py <args> # Run script
```

## Connection and Authentication

Connection details for your environments should be specified in the format described in [`config.json.tmpl`](config.json.tmpl).

Each environment must be distinctly named.

## Running

Full usage information:

```
usage: contrast_application_licenses.py [-h] [-c CONFIG_FILE] [-i UPDATE_INTERVAL] [-l {CRITICAL,ERROR,WARN,INFO,DEBUG}] [-p PROMETHEUS_LISTEN_PORT | -u PROMETHEUS_PUSH_GATEWAY]

Utility to count licensed Contrast Assess Applications across environments, de-duplicating them by name, language, and metadata.

options:
  -h, --help            show this help message and exit
  -c CONFIG_FILE, --config_file CONFIG_FILE, --config-file CONFIG_FILE
                        Path to JSON config or - to read it from stdin, defaults to config.json
  -i UPDATE_INTERVAL, --update-interval UPDATE_INTERVAL, --update_interval UPDATE_INTERVAL
                        Number of minutes to wait between polls of the configured environments for licensed applications. Only used when serving prometheus data with -p.
  -l {CRITICAL,ERROR,WARN,INFO,DEBUG}, --log-level {CRITICAL,ERROR,WARN,INFO,DEBUG}, --log_level {CRITICAL,ERROR,WARN,INFO,DEBUG}
                        Log level
  -p PROMETHEUS_LISTEN_PORT, --prometheus-listen-port PROMETHEUS_LISTEN_PORT, --prometheus_listen_port PROMETHEUS_LISTEN_PORT
                        Port to serve metrics on.
  -u PROMETHEUS_PUSH_GATEWAY, --prometheus-push-gateway PROMETHEUS_PUSH_GATEWAY, --prometheus_push_gateway PROMETHEUS_PUSH_GATEWAY
                        URL for a Prometheus push gateway where metrics will be sent.
```

If used with `-p`, the license data will be periodically refreshed (default every 5 minutes), and served on the specified port (daemon mode).

If used with `-u`, the license data is retrieved once and sent to the specified push gateway URL. This is good for cron-style environments.

Both options may not be used together.

If neither option is provided, counts are logged at the default info level.

## Output

```
# HELP contrast_assess_unique_licensed_applications Number of unique licensed Contrast Assess applications, de-duplicated by name, language and metadata values.
# TYPE contrast_assess_unique_licensed_applications gauge
contrast_assess_unique_licensed_applications 4.0
# HELP contrast_assess_licensed_applications Number of licensed Contrast Assess applications on an environment.
# TYPE contrast_assess_licensed_applications gauge
contrast_assess_licensed_applications{environment="Environment1"} 3.0
contrast_assess_licensed_applications{environment="Environment2-EU"} 3.0

```

## Development Setup
Various tools enforce code standards, and are run as a pre-commit hook. This must be setup before committing changes with the following commands:
```bash
python3 -m venv venv # setup a virtual environment
. venv/bin/activate # activate the virtual environment
pip3 install -r requirements-dev.txt # install development dependencies (will also include app dependencies)
pre-commit install # setup the pre-commit hook which handles formatting
```
