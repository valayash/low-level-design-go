# Snake and Ladder

**Difficulty:** Easy–Medium · **Focus:** finding the right abstraction, dependency injection, testing randomness

Tic Tac Toe tested whether you could *separate* responsibilities. This one tests
something harder: whether you can spot that two things you think are different
are actually the same thing — and whether you can test code that depends on
randomness.

---

## Requirements

**Functional**
1. A board of `N` cells, numbered `1..N` (classically 100). Players start off-board at 0.
2. Two or more players take turns rolling a die (1–6) and advancing by the roll.
3. Landing on the **head of a snake** sends the player down to its tail.
4. Landing on the **bottom of a ladder** sends the player up to its top.
5. The first player to reach cell `N` **exactly** wins.
6. A roll that would overshoot `N` is forfeited — the player stays put.

**Non-functional**
- No `fmt.Println` in the game logic. Same rule as last time.
- **Every test must be deterministic.** A test that sometimes passes is not a test.
- A board that is impossible to win must be rejected at construction, not discovered at runtime.

---

## Design first

### The question this problem exists to ask

You are about to write a `Snake` type and a `Ladder` type. Before you do — write down
what each one *stores* and what each one *does*.

A snake: you land on cell 17, you end up at cell 4.
A ladder: you land on cell 9, you end up at cell 31.

**Are those two different behaviours, or one behaviour with different numbers?**

If you build two types, you'll write the same field pair twice, the same lookup twice,
and every future feature (a teleporter? a trapdoor?) is a third copy. If you build one,
what do you call it, and how does anything tell a snake from a ladder when it needs to —
say, to render them differently?

There's a real argument on both sides. Pick one and be ready to defend it, because
this is the question the interviewer is actually asking.

### The question that decides whether your code is testable

Your game needs a die. The obvious move:

```go
roll := rand.Intn(6) + 1
```

Now write a test proving that a player landing on a snake head slides to the tail.

You can't. You cannot make the die produce the number your test needs. Your only
options are to run the game a thousand times and hope, or to give up and test nothing.

**So: how does the die get into your `Game`?** What's the type of the thing that
produces a roll, and who supplies it? Once you answer that, ask what your test
passes in instead.

This is dependency injection, and it's the single most valuable idea in this problem.

### The rest

1. **How do you store the jumps?** A `[]Snake` and a `[]Ladder` you scan? A `map[int]int`?
   What's the lookup cost of each, and which one makes "two snakes with the same head"
   impossible to represent?

2. **Where does a player's position live?** On `Player`, or in the `Game`? You settled
   this exact argument for turn state in Tic Tac Toe — does the same reasoning apply?

3. **What must a constructor reject?** A snake whose tail is above its head. A ladder
   from cell 1 (nobody can land there... or can they?). A ladder ending on a snake head.
   Two jumps starting from the same cell. A snake on the final cell. Which of these are
   genuinely broken, and which are just unusual?

4. **What does one turn actually do?** Write the steps in order, like `Move` last time.
   There are more than you think — and one of them can chain.

5. **Can a jump land you on another jump?** Ladder to cell 31, and 31 is a snake head.
   Do you slide again? Real rules say yes. What does that do to your code, and what
   stops an infinite loop?

---

## Follow-ups

- Two dice instead of one.
- Rolling a 6 grants another turn. Three 6s in a row forfeits it.
- A "crocodile" that moves you *sideways* — how much do you rewrite?
- Emit events (`playerMoved`, `snakeBite`, `gameWon`) for a UI to subscribe to. Which pattern?
- Make the game safe under concurrent moves.

---

## Status

- [ ] Design sketched
- [ ] Jump abstraction decided
- [ ] Die injected, not hardcoded
- [ ] Board construction validated
- [ ] Turn logic
- [ ] Tests passing (`go test ./problems/snakeandladder/ -v`)
