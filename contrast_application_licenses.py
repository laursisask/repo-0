import argparse
import logging
import sys
from collections import defaultdict
from dataclasses import dataclass
from time import sleep

from contrast_api import contrast_instance_from_json, load_config

args_parser = argparse.ArgumentParser(
    description="Utility to count licensed Contrast Assess Applications across environments, de-duplicating them by name, language, and metadata."
)
args_parser.add_argument(
    "-c",
    "--config_file",
    "--config-file",
    help="Path to JSON config or - to read it from stdin, defaults to config.json",
    default="config.json",
    type=argparse.FileType("r"),
)
args_parser.add_argument(
    "-i",
    "--update-interval",
    "--update_interval",
    help="Number of minutes to wait between polls of the configured environments for licensed applications. Only used when serving prometheus data with -p.",
    type=int,
    default=5,
)
args_parser.add_argument(
    "-l",
    "--log-level",
    "--log_level",
    help="Log level",
    choices=["CRITICAL", "ERROR", "WARN", "INFO", "DEBUG"],
    type=str.upper,
    default="INFO",
)
group = args_parser.add_mutually_exclusive_group()
group.add_argument(
    "-p",
    "--prometheus-listen-port",
    "--prometheus_listen_port",
    help="Port to serve metrics on.",
    type=int,
)
group.add_argument(
    "-u",
    "--prometheus-push-gateway",
    "--prometheus_push_gateway",
    help="URL for a Prometheus push gateway where metrics will be sent.",
)
args = args_parser.parse_args()

logging.basicConfig(level=args.log_level, format="%(levelname)s: %(message)s")
logger = logging.getLogger(__file__)

try:
    from prometheus_client import (
        CollectorRegistry,
        Gauge,
        push_to_gateway,
        start_http_server,
    )
except ImportError:
    logger.fatal("prometheus-client is not installed")
    sys.exit(1)

try:
    import schedule
except ImportError:
    logger.fatal("schedule is not installed")
    sys.exit(1)


config = load_config(file=args.config_file)

environments = {}
for iteration, org in enumerate(config):
    environment = contrast_instance_from_json(org)
    name = org["name"]

    if name in environments:
        logger.error(
            f"Environment named '{name}' was already added prior to environment[{iteration}], please use distinct names"
        )
        sys.exit(1)

    if not environment.test_connection():
        logger.error(f"Test connection failed for Environment '{name}'")
        sys.exit(1)
    if not environment.test_org_access():
        logger.error(f"Organization access failed for environment '{name}'")
        sys.exit(1)

    environments[name] = environment


@dataclass(eq=True, frozen=True)
class Application:
    """Dataclass to represent an application using name, language and metadata as composite unique key."""

    name: str
    language: str
    metadata: str


def metadata_to_str(metadata_entities: dict) -> str:
    """Convert metadata dictionary from the application response into a string like the YAML format to provide application metadata."""
    return ",".join(
        map(
            lambda meta: f"{meta['fieldName']}={meta['fieldValue']}",
            sorted(metadata_entities, key=lambda meta: meta["fieldName"]),
        )
    )


@dataclass
class AppCounts:
    unique_licensed_applications: int
    environment_applications: dict[str, dict]


def count_licensed_applications() -> AppCounts:
    apps: set[Application] = set()
    # map of environment -> language -> application count
    environment_apps_by_language: dict[str, dict[str, int]] = {}

    for environment_name, environment in environments.items():
        # map of language -> application count, defaulting at 0 for new keys
        language_count = defaultdict(int)
        logger.info(f"Listing applications for environment '{environment_name}'...")
        applications = environment.list_org_apps(
            environment._org_uuid,
            include_archived=True,
            include_merged=False,
            quick_filter="LICENSED",
        )
        for application in applications:
            metadata = metadata_to_str(application["metadataEntities"])

            app = Application(
                application["name"],
                application["language"],
                metadata,
            )
            apps.add(app)
            language_count[app.language] = language_count[app.language] + 1

        logger.debug(f"Unique application count is now: {len(apps)}")

        environment_app_count = len(applications)
        logger.info(
            f"Environment '{environment_name}' app count: {environment_app_count}"
        )
        environment_apps_by_language[environment_name] = language_count

    unique_license_count = len(apps)
    logger.info(f"Unique license count: {unique_license_count}")
    return AppCounts(unique_license_count, environment_apps_by_language)


registry = CollectorRegistry()
unique_guage = Gauge(
    "contrast_assess_unique_licensed_applications",
    "Number of unique licensed Contrast Assess applications, de-duplicated by name, language and metadata values.",
    registry=registry,
)
environment_guage = Gauge(
    "contrast_assess_licensed_applications_total",
    "Number of licensed Contrast Assess applications on an environment.",
    ["environment"],
    registry=registry,
)
language_guage = Gauge(
    "contrast_assess_licensed_applications",
    "Number of licensed Contrast Assess applications in a specific language.",
    ["language", "environment"],
    registry=registry,
)


def update_registry():
    data = count_licensed_applications()
    unique_guage.set(data.unique_licensed_applications)

    for environment, language_counts in data.environment_applications.items():
        total_apps_in_environment = 0

        for language, count in language_counts.items():
            total_apps_in_environment += count
            gauge = language_guage.labels(environment=environment, language=language)
            gauge.set(count)

        gauge = environment_guage.labels(environment=environment)
        gauge.set(total_apps_in_environment)


update_registry()

if listen_port := args.prometheus_listen_port:
    start_http_server(listen_port, registry=registry)
    logger.info(f"Listening on port {listen_port}")
    logger.info(f"Scheduling update every {args.update_interval} minute(s)")
    schedule.every(args.update_interval).minutes.do(update_registry)
    while True:
        schedule.run_pending()
        sleep(1)

if url := args.prometheus_push_gateway:
    logger.info(f"Pushing data to gateway at {url}...")
    push_to_gateway(url, job="contrast_assess_licenses_used", registry=registry)
    logger.info("Successfully pushed data to gateway.")
