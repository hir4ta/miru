"""sample.py — exercises common Python syntax features."""

from __future__ import annotations

import asyncio
from dataclasses import dataclass
from typing import Iterable

PI: float = 3.14159


def greet(who: str = "world") -> str:
    """Return a friendly greeting."""
    return f"hello, {who}"


@dataclass(frozen=True)
class Circle:
    radius: float

    @property
    def area(self) -> float:
        return PI * self.radius**2


async def fetch_user(client, id: int) -> dict | None:
    try:
        res = await client.get(f"/users/{id}")
        res.raise_for_status()
        return res.json()
    except Exception as err:
        print(f"fetch failed: {err}")
        return None


def summarize(items: Iterable[str], *, tag: str = "default") -> dict[str, object]:
    first, *rest = list(items) or [None]
    return {"tag": tag, "first": first, "rest_count": len(rest)}


# List comprehension + generator + walrus
words = ["alpha", "beta", "gamma", "delta"]
a_words = [w for w in words if w.startswith("a")]
print(a_words)

if (n := len(words)) > 2:
    print(f"got {n} words")


if __name__ == "__main__":
    print(greet("miru"))
    print(f"area: {Circle(5).area}")
    asyncio.run(asyncio.sleep(0))
