# Online Shopping Service — Design Notes

The current model. Edited in place as decisions change — this always reflects where the
design stands now, not how it got here. No Go in this file; types and relationships only.

---

## Entities

![entity diagram](design.svg)

| Entity | Holds | Purpose |
|---|---|---|
| `User` | ID, name, email, their orders | Who is shopping, and their history |
| `Order` | Order ID, line items, total | A basket at checkout and what it came to |
| `ProductItem` | A product, a quantity | One line: *"2 of this one"* |
| `Product` | ID, name, description, stock, price | A thing that can be bought, and how many remain |

## Decisions

- **Order history hangs off the user** — `User.Orders` is that user's orders.
- **Line items, not bare products** — `Order.Items` is `[]ProductItem`, so quantity has a home.
- **The order is the cart** — no separate `Cart` entity; the order's item list plays both roles.
- **Stock lives on the product** — `Product.Inventory` is the count remaining.
- **The order stores its own total** — `CartValue` is carried, not recomputed.

---

## Open questions

Each of these changes the diagram above.

### 1. Price at the time of purchase

`Price` lives on `Product`, and `ProductItem` points at a live product. Buy a camera for
£500 in January; drop its price to £400 in March; your January order now reports £400. The
financial record changes retroactively.

**Q:** what must an order line *copy* rather than reference — and does `Order` get to point
at `Product` at all?

### 2. Cart line vs order line

A cart line *should* track the live price: if the camera drops while it sits in the basket,
you want the lower price. An order line must be frozen forever. Same two fields, opposite
requirements.

**Q:** one type used two ways, or two types?

### 3. `CartValue` is derived

The total is the sum of the lines. Storing it means every mutation must update it, and any
path that forgets leaves the order lying about its own value.

**Q:** stored or computed? (For an *order*, unlike a cart, there is a real argument for
storing it — see question 1.)

### 4. Nothing tracks order status

The requirements need `Pending → Processing → Shipped → Delivered`, plus `Cancelled`.

**Q:** what field, what type, and what stops `Delivered → Pending`?

### 5. `User.Orders` makes a cycle

If an order ever needs to know whose it is, `User → Order → User` is a reference cycle. It
also means loading one user drags in their whole order history.

**Q:** does the user hold orders, or does the order hold a user ID with the lookup elsewhere?

### 6. Stock on the product is the race

Two users check out the last unit at once. Both read `Inventory == 1`, both decrement,
stock goes to `-1`.

**Q:** what single operation is atomic, and when is stock actually taken — add-to-cart,
checkout, or payment success? What returns it if payment then fails?

### 7. The type of `money`

`float64` cannot represent `0.1 + 0.2` exactly, so totals drift and equality breaks.

**Q:** what type, and in what unit?

### 8. Payment

Multiple payment methods are required, and a payment can fail in ways the caller must tell
apart — declined, insufficient funds, timeout.

**Q:** what is the interface, its one method, and what does that method return?
