"""Write Synapse's config from the environment, then start Synapse.

The image carries templates only. Every secret arrives through the
environment, so nothing sensitive is baked into the image or the repository.
"""

import os
import re
import subprocess
import sys

SYNAPSE_UID = 991

REQUIRED = [
    "MATRIX_SERVER_NAME",
    "MATRIX_PUBLIC_URL",
    "SYNAPSE_DB_HOST",
    "SYNAPSE_DB_USER",
    "SYNAPSE_DB_PASSWORD",
    "SYNAPSE_MACAROON_SECRET",
    "SYNAPSE_FORM_SECRET",
    "PRISME_AS_TOKEN",
    "PRISME_HS_TOKEN",
]


def render(source, target, values):
    with open(source) as handle:
        text = handle.read()
    text = re.sub(r"\$\{([A-Z0-9_]+)\}", lambda match: values[match.group(1)], text)
    with open(target, "w") as handle:
        handle.write(text)


def write_signing_key(key):
    """Write the configured key, or generate one once when none is given."""
    if key:
        with open("/data/signing.key", "w") as handle:
            handle.write(key + "\n")
        return
    if not os.path.exists("/data/signing.key"):
        subprocess.run(
            [sys.executable, "-m", "synapse._scripts.generate_signing_key", "-o", "/data/signing.key"],
            check=True,
        )


def main():
    missing = [name for name in REQUIRED if not os.environ.get(name)]
    if missing:
        sys.exit("synapse: missing environment: " + ", ".join(missing))
    values = dict(os.environ)
    values.setdefault("SYNAPSE_DB_NAME", "synapse")
    values["MATRIX_SERVER_NAME_REGEX"] = re.escape(values["MATRIX_SERVER_NAME"]).replace("\\", "\\\\")
    values["MATRIX_PUBLIC_URL"] = values["MATRIX_PUBLIC_URL"].rstrip("/")
    os.makedirs("/data/media_store", exist_ok=True)
    for name in ("homeserver.yaml", "appservice.yaml", "log.config"):
        render("/prisme/" + name, "/data/" + name, values)
    write_signing_key(values.get("SYNAPSE_SIGNING_KEY", "").strip())
    os.chmod("/data/signing.key", 0o600)
    for path in ("/data", "/data/media_store", "/data/signing.key"):
        os.chown(path, SYNAPSE_UID, SYNAPSE_UID)
    os.environ["SYNAPSE_CONFIG_PATH"] = "/data/homeserver.yaml"
    os.execv("/start.py", ["/start.py"])


main()
