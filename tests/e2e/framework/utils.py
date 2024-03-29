import logging
import random
import string

LOG = logging.getLogger(__name__)


def random_name(prefix: str = "", N: int = 6) -> str:
    return prefix + "".join(
        random.choices(string.ascii_lowercase, k=1)
        + random.choices(string.ascii_lowercase + string.digits, k=N - 1)
    )
