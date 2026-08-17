# Splitwise

**Difficulty:** Medium · **Focus:** Strategy pattern, modelling money, invariants that must never break

The first genuinely *hard* problem in this list. Tic Tac Toe was about separating
responsibilities; Snake and Ladder was about injecting a dependency. This one is
about **modelling a domain where being wrong by one cent is a bug**, and where the
"obvious" data type for money is the wrong one.

---

## Requirements

**Functional**
1. Users can be registered. Each has an id, a name, and contact details.
2. A user can add an **expense**: an amount, who paid, and who it's split between.
3. Splits come in three flavours:
   - **Equal** — divide evenly among participants
   - **Exact** — caller names each participant's exact share
   - **Percent** — caller names each participant's percentage
4. The system tracks **who owes whom, and how much**.
5. A user can see their balances: total owed to them, total they owe, and the per-person breakdown.
6. Users can **settle up** — record a payment that reduces a debt.
7. Optional: **groups**, so expenses can be scoped to a set of people.

**Non-functional**
- The sum of all balances in the system must always be **exactly zero**. If it drifts, money was invented or destroyed.
- Splits that don't add up must be rejected at entry, never stored.
- No `fmt.Println` in the domain logic.
- Deterministic tests. Money bugs hide in the cases you didn't write down.

---

## Design first

### The question that fails most candidates

You're about to write this:

```go
type Expense struct {
	Amount float64
}
```

**Don't.** Try this in Go and look at the output:

```go
fmt.Println(0.1 + 0.2)           // 0.30000000000000004
fmt.Println(0.1+0.2 == 0.3)      // false
```

`float64` cannot represent most decimal fractions exactly. Every split, every
addition, every comparison accumulates error. Balances that should cancel to zero
won't, and `if balance == 0` will be false for a settled account.

So: **what type do you use for money, and what unit does it hold?** Once you answer
that, a second question follows immediately — what does your API take from a caller,
and where does the conversion happen?

### The question about the split itself

Three split types today. Tomorrow the interviewer adds "split by shares" (Alice gets
2 shares, Bob gets 1). Write down what your code looks like if you handle this with:

```go
switch splitType {
case Equal:   ...
case Exact:   ...
case Percent: ...
}
```

Now add the fourth type. How many places did you edit? What if two of those `switch`
statements exist in different files and you only remembered one?

**What's the alternative?** You already met it — it's the same shape as the
`WinningStrategy` in Tic Tac Toe and the injected die in Snake and Ladder. Name the
interface, give it one method, and say what that method takes and returns.

### The question nobody asks until it bites

Three people split £10 equally. Each owes £3.333...

In your chosen money type, what are the three shares, and **do they sum to the
original amount?** If they don't, you just destroyed a penny — and your
"balances must sum to zero" invariant is broken on the very first expense.

Who eats the remainder, and is that decision deterministic? (If the answer is "the
first person in the map", be aware Go randomises map iteration order, so the same
expense would round differently on different runs.)

### The rest

1. **How do you store balances?** Two candidates:
   - A pairwise ledger: `balances[alice][bob] = 500` means Alice owes Bob £5
   - A single net figure per user: `net[alice] = -500`

   One of these can answer "who exactly do I owe?" and the other can't. One makes
   "simplify the group's debts" easy and the other doesn't. Which is which?

2. **What validates a split, and when?** Exact shares must sum to the total.
   Percentages must sum to 100. A participant list must be non-empty and free of
   duplicates. Is that the `Expense` constructor's job, the split strategy's, or
   the service's?

3. **Can you pay yourself?** What happens if the payer is also in the participant
   list — which is the *normal* case? Work through the arithmetic.

4. **What's the public API?** Write the signatures a caller uses. Something like
   `AddExpense(...)`, `Settle(...)`, `Balances(userID)`. Be precise about what each
   returns and how each fails.

5. **Where does settling up fit?** Is a settlement just an expense with a special
   shape, or its own concept? Argue it either way.

---

## Follow-ups

- **Simplify debts**: Alice owes Bob £10, Bob owes Carol £10 → Alice owes Carol £10, and Bob is clear. What's the algorithm, and what's its complexity?
- Multiple currencies.
- Expense history and an audit trail — can you recompute every balance from the log?
- Notify users when they're added to an expense. Which pattern?
- Concurrent expense entry from several users.

---

## Status

- [ ] Design sketched
- [ ] Money representation decided
- [ ] Split strategy interface
- [ ] Balance storage
- [ ] Validation on entry
- [ ] Settle up
- [ ] Tests passing (`go test ./problems/splitwise/ -v`)
- [ ] Follow-up: simplify debts
