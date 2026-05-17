// sample.js — exercises common JavaScript syntax features.

const PI = 3.14159;
let name = "miru";

// Template literal + arrow function
const greet = (who) => `hello, ${who}`;

console.log(greet(name));

// Class with private field + getter
class Circle {
  #radius;
  constructor(radius) {
    this.#radius = radius;
  }
  get area() {
    return PI * this.#radius ** 2;
  }
}

const c = new Circle(5);
console.log(`area: ${c.area}`);

// Async / await + try-catch
async function fetchUser(id) {
  try {
    const res = await fetch(`/users/${id}`);
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    return await res.json();
  } catch (err) {
    console.error("fetch failed:", err);
    return null;
  }
}

// Destructuring + spread + default params
function summarize({ items = [], tag = "default" }) {
  const [first, ...rest] = items;
  return { tag, first, restCount: rest.length };
}

// Regex + array methods
const lines = "alpha\nbeta\ngamma".split("\n").filter((l) => /^a/.test(l));
console.log(lines);

export { Circle, greet, fetchUser, summarize };
