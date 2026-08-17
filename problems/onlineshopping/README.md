# Online Shopping Service

**Difficulty:** Hard · **Focus:** state machines, inventory under concurrency, several patterns at once

The biggest problem so far. Tic Tac Toe was one type doing one job; this is a dozen types
that have to agree with each other. The interviewer is no longer asking *"can you separate
concerns"* — they're asking *"can you keep a system consistent when two people race for the
last item in stock?"*

Scope discipline matters here. A real Amazon is a hundred services. You have 45 minutes.
Decide what you're **not** building and say so out loud.

---

## Requirements

**Functional**
1. Browse products by category, and search them by name.
2. Add products to a cart, change quantities, remove them.
3. Place an order from a cart. Pay for it.
4. Track order status through its lifecycle; cancel while that's still legal.
5. Manage a user profile and view order history.
6. Inventory must reflect reality — you cannot sell what you do not have.
7. Support more than one payment method.

**Non-functional**
- **Two users racing for the last unit must not both succeed.** This is the requirement that makes the problem hard.
- Money is not `float64`. Same rule as Splitwise.
- An order's status may only follow legal transitions. `Delivered → Pending` must be impossible.
- No `fmt.Println` in the domain logic.
- Deterministic tests, including for the race.

---

## Design first

### The question this problem exists to ask

Two customers have the same last-in-stock camera in their carts. Both hit checkout at the
same instant.

```go
if product.Quantity >= requested {   // both read 1, both pass
	product.Quantity -= requested    // both decrement
}
```

You've now sold two cameras and have `-1` in stock.

- **Where does the stock live**, and what guards it?
- **When is stock actually taken** — when it goes in the cart, at checkout, or when payment clears? Each answer has a cost. Reserving at cart time means an abandoned cart holds inventory hostage; reserving at payment time means a customer can pay for something that just sold out.
- **What's the smallest thing you can lock**, and what happens to throughput if you get that wrong?
- If payment fails after you've decremented stock, **what puts it back?**

Write down your answer before you write any code. This is the whole problem.

### The question about order status

`Pending → Processing → Shipped → Delivered`, plus `Cancelled`. The obvious version:

```go
order.Status = Shipped   // any code, anywhere, any time
```

Nothing stops `Delivered → Pending`, or cancelling something already delivered, or shipping
an unpaid order.

- **Who is allowed to change status, and what enforces the legal transitions?**
- Where does the transition table live — a `map`, a `switch`, or a type per state?
- Cancellation is legal in some states and not others. Where does that rule go so that it exists in exactly one place?

### The rest

1. **Cart versus Order.** A cart holds a product and a quantity. An order holds a product and
   a quantity. Are they the same type? Think hard before saying yes: what happens to an order
   from last year when the product's price changes today? Which one must remember the price
   *at the time*, and what does that tell you?

2. **Payment.** You met this shape in Splitwise's splits and Tic Tac Toe's win rules. Name the
   interface and its one method. Then: what does it return, given that a payment can fail for
   reasons the caller must distinguish (declined, insufficient funds, network timeout)?

3. **Search.** "Search by name" is easy to fake with a linear scan. State the complexity of
   your approach and what you'd change at a million products — an interviewer will ask, and
   "I'd add an index" with a reason beats a clever answer you can't justify.

4. **The service object.** The reference design makes this a Singleton. Do you agree? In Go, a
   package-level instance and a `sync.Once` is the usual way — but a Singleton is also global
   mutable state that makes tests interfere with each other. Argue it either way, and say how
   you'd test it.

5. **What's the public API?** Signatures for `AddToCart`, `PlaceOrder`, `Cancel`, `Search`.
   Be precise about failure: out of stock, payment declined, illegal transition, and product
   not found are four different problems.

---

## Scope: say what you are not building

Out of scope unless asked — and say so rather than letting the interviewer wonder:
authentication, a real payment gateway, shipping/logistics, reviews, recommendations,
persistence, HTTP handlers.

---

## Follow-ups

- Discount codes and promotions. Where do they apply without touching `Order`?
- Notify a user when order status changes. Which pattern?
- Multiple warehouses — stock is per-location now.
- Partial shipment: one order, two parcels.
- Returns and refunds. What does that do to your state machine?

---

## Status

- [ ] Design sketched
- [ ] Money representation decided
- [ ] Inventory + concurrency strategy decided
- [ ] Order state machine
- [ ] Payment interface
- [ ] Cart and checkout
- [ ] Tests passing, including a `-race` test for the last-item race
- [ ] Follow-up: discounts
