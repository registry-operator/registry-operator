import kr8s
from kr8s.objects import Namespace

from e2e.framework.utils import random_name, LOG


def kubernetes_cluster() -> None:
    kr8s.whoami()


def namespace_create() -> str:
    ns = Namespace(
        {
            "apiVersion": "v1",
            "kind": "Namespace",
            "metadata": {
                "name": random_name(prefix="e2e-"),
            },
        }
    )

    ns.create()

    return ns.name


def namespace_delete(name: str) -> None:
    ns = Namespace(
        {
            "metadata": {
                "name": name,
            }
        }
    )

    ns.delete()


def namespace_exists(name: str) -> None:
    LOG.info(f'checking if namespace "{name}" exists')
    Namespace.get(name=name)
