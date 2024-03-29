import time
from typing import Generator

import pytest
from pytest_bdd import given, scenario, then

import e2e.framework as framework


@pytest.fixture
def namespace() -> Generator[str, None, None]:
    ns = framework.namespace_create()
    yield ns
    framework.namespace_delete(name=ns)


@scenario("features/sanity.feature", "Testing fixture")
def test_simplest_registry() -> None:
    pass


@given("kubernetes cluster")
def kubernetes_cluster() -> None:
    framework.kubernetes_cluster()


@then("e2e namespace should exist")
def namespace_exists(namespace) -> None:
    framework.namespace_exists(name=namespace)
