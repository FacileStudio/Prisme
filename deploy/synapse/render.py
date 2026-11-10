"""Write Synapse's config from the environment, then start Synapse.

The image carries templates only. Every secret arrives through the
environment, so nothing sensitive is baked into the image or the repository.
"""

import json
import os
import re
import subprocess
import sys
import time
import urllib.error
import urllib.request

SYNAPSE_UID = 991
SENDER_LOCALPART = "prisme"
LOCAL_URL = "http://127.0.0.1:8008"

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


def set_admin(values, sender):
    """Flag the sender as a server admin. False when it has no row yet."""
    import psycopg2

    connection = psycopg2.connect(
        host=values["SYNAPSE_DB_HOST"],
        user=values["SYNAPSE_DB_USER"],
        password=values["SYNAPSE_DB_PASSWORD"],
        dbname=values["SYNAPSE_DB_NAME"],
    )
    try:
        with connection, connection.cursor() as cursor:
            cursor.execute("SELECT to_regclass('users')")
            if cursor.fetchone()[0] is None:
                return False
            cursor.execute("UPDATE users SET admin = 1 WHERE name = %s", (sender,))
            return cursor.rowcount == 1
    finally:
        connection.close()


def wait_healthy(synapse):
    for _ in range(120):
        if synapse.poll() is not None:
            raise RuntimeError("synapse exited during its first start")
        try:
            with urllib.request.urlopen(LOCAL_URL + "/health", timeout=2):
                return
        except OSError:
            time.sleep(1)
    raise RuntimeError("synapse did not become healthy")


def register_sender(values):
    body = {"type": "m.login.application_service", "username": SENDER_LOCALPART, "inhibit_login": True}
    request = urllib.request.Request(
        LOCAL_URL + "/_matrix/client/v3/register",
        data=json.dumps(body).encode(),
        headers={"Authorization": "Bearer " + values["PRISME_AS_TOKEN"], "Content-Type": "application/json"},
    )
    try:
        urllib.request.urlopen(request, timeout=10).close()
    except urllib.error.HTTPError as error:
        if error.code != 400:
            raise


def register_once(values):
    """Run Synapse once, so the sender gets a row, then stop it."""
    synapse = subprocess.Popen(["/start.py"])
    try:
        wait_healthy(synapse)
        register_sender(values)
    finally:
        synapse.terminate()
        synapse.wait(timeout=60)


def promote_sender(values):
    """Make the appservice sender a server admin before Synapse starts.

    The API calls Synapse's admin API with the appservice token, which passes
    only while the sender is an admin. Synapse caches that flag per process,
    so it is set while Synapse is down. On a first start the sender has no row
    yet, so Synapse runs once to register it. A failure is printed and Synapse
    still starts: chat never waits on the admin API.
    """
    sender = "@" + SENDER_LOCALPART + ":" + values["MATRIX_SERVER_NAME"]
    try:
        if set_admin(values, sender):
            return
        register_once(values)
        if not set_admin(values, sender):
            print("synapse: the appservice sender has no row to promote", file=sys.stderr)
    except Exception as error:
        print("synapse: could not make the appservice sender an admin: %s" % error, file=sys.stderr)


def main():
    missing = [name for name in REQUIRED if not os.environ.get(name)]
    if missing:
        sys.exit("synapse: missing environment: " + ", ".join(missing))
    values = dict(os.environ)
    values.setdefault("SYNAPSE_DB_NAME", "synapse")
    values["PRISME_AS_URL"] = values.get("PRISME_AS_URL") or "null"
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
    promote_sender(values)
    os.execv("/start.py", ["/start.py"])


main()
