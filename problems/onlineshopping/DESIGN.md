# Online Shopping Service — Design Notes

The current model. Edited in place as decisions change — this always reflects where the
design stands now, not how it got here. No Go in this file; types and relationships only.

Structure aligned with the reference solution in
[awesome-low-level-design](https://github.com/ashishps1/awesome-low-level-design/tree/main/solutions/golang/onlineshoppingservice),
with the deviations listed at the bottom.

---

## Entities

![entity diagram](design.svg)

| Type | Kind | Holds | Purpose |
|---|---|---|---|
| `OnlineShoppingService` | service | `Users`, `Products`, `Orders` maps | Owns the registries; runs checkout |
| `User` | entity | ID, name, email, password, orders | Who is shopping |
| `Product` | entity | ID, name, description, price, quantity | A sellable thing and its stock |
| `ShoppingCart` | entity | `items map[productID]OrderItem` | The mutable basket before checkout |
| `OrderItem` | value | product, quantity | One line: *"2 of this one"* |
| `Order` | entity | ID, user, items, total, status | The record of a completed purchase |
| `OrderStatus` | enum | Pending, Processing, Shipped, Delivered, Cancelled | Where an order is in its lifecycle |
| `Payment` | interface | `ProcessPayment(amount)` | Swappable payment method |
| `CreditCardPayment` | impl | — | One concrete method |

## Layers

- **Service** — verbs that span entities. `PlaceOrder` touches inventory, orders, cart and
  payment, so it belongs to none of them.
- **Entities** — nouns that own only rules about themselves. `Product` knows whether it has
  stock; `ShoppingCart` knows how to merge a duplicate line.

The three maps stand in for a database: `s.Products[id] = p` is an `INSERT`, and
`s.Products[id]` is a `SELECT ... WHERE id = ?`.

## Decisions

- **One line type, `OrderItem`, shared by cart and order.** The cart stores it keyed by
  product ID so duplicates merge rather than creating two lines.
- **`Order` is created by a constructor** that computes the total and sets `Pending`. It is
  never assembled field by field.
- **Status is an enum**, not a string, so an invalid status cannot be typed.
- **Payment is an interface** with one method, so a new method is a new type rather than a
  new `switch` branch.
- **Service is a singleton**, holding the only copies of users, products and orders.
- **Order history is stored twice** — on `User.Orders` and in `Service.Orders`.

## What the checkout does, in order

```
for each cart line:  if stock available  →  deduct it
if nothing available →  error
create the Order (status Pending), store it, append to user, clear the cart
take payment
  success →  status Processing
  failure →  status Cancelled, put all the stock back
```

---

## Open questions

These are places the reference makes a choice worth arguing with. Each is still undecided
for our version.

### 1. Money is `float64` in the reference

`0.1 + 0.2 != 0.3`. Totals accumulate error, and equality comparisons on money stop working.

**Q:** integer minor units instead — and if so, what unit, and where does conversion happen?

### 2. The freeze is not real

`OrderItem` holds `*Product`, so a line's price is read live from the catalogue. `TotalAmount`
*is* snapshotted at construction, so the total is stable — but the line items it was
computed from are not. Change a price and an order's total no longer equals the sum of its
own lines.

**Q:** does `OrderItem` copy `ProductID` and `UnitPrice` instead of holding a pointer?

### 3. Order history lives in two places

`user.AddOrder(*order)` appends a **copy** of the order, while `s.Orders[id]` keeps the
pointer. `SetStatus` afterwards mutates the pointer only, so the user's own history keeps
whatever status it was copied with.

**Q:** one home for orders — the registry, with `OrdersFor(userID)` as a query?

### 4. The stock race is unguarded

`IsAvailable(n)` then `UpdateQuantity(-n)` is check-then-act with no lock. Two checkouts
both see the last unit and both take it.

**Q:** what single operation is atomic, and who holds the mutex?

### 5. Partial orders happen silently

Unavailable lines are skipped, and the order is placed with whatever was in stock. The
customer is charged for less than they asked for and is never told.

**Q:** all-or-nothing, or partial with an explicit report of what was dropped?

### 6. Payment runs after the cart is cleared and the order is stored

If payment fails, the order is already in the registry and the basket is already empty. The
customer loses their cart on a declined card.

**Q:** what is the correct order of operations, and what happens if the process dies between
deducting stock and refunding it?

### 7. `ProcessPayment` returns `bool`

Declined, insufficient funds, and a network timeout are different problems — one is the
customer's, one is retryable, one is neither.

**Q:** does it return an `error` so callers can tell them apart?

### 8. The singleton is not thread-safe

`if instance == nil` is itself a race, in the type responsible for handling concurrency.
It also makes every test share one catalogue.

**Q:** `sync.Once`, or inject the service so tests get a fresh one?

### 9. Order IDs collide

The reference builds them as `"ORDER-" + user.ID`, so a user's second order overwrites their
first in the registry.

**Q:** what generates a unique order ID?
