# Tic Tac Toe

**Difficulty:** Easy · **Focus:** enums, interfaces, encapsulation, state management

The classic LLD warm-up. It looks trivial, but interviewers use it to see whether you
can separate *rules* from *state* from *I/O* — and whether your design survives the
follow-up questions at the bottom of this file.

---

## Requirements

**Functional**
1. Two players alternate turns on an `N x N` board (default `3 x 3`).
2. Each player has a distinct symbol (`X` and `O`).
3. A player marks one empty cell per turn.
4. The game ends when a player fills a full row, column, or either diagonal — that player wins.
5. If the board fills with no winner, the game is a **draw**.
6. Invalid moves must be rejected: out of bounds, cell already taken, game already over,
   or a player moving out of turn.

**Non-functional**
- The core game logic must have **no `fmt.Println`** in it. Rules and rendering are separate concerns.
- Illegal states should be unrepresentable where the type system allows it.
- Every rule above should be provable by a test.

---

## Design first — before you write any code

LLD is not "start typing structs." Answer these in writing, then we build it.

1. **What are your types?** Name every struct and interface you think you need, and the
   one responsibility each carries. If a type has two responsibilities, split it.

2. **How do you represent a cell?** A `string`? A `rune`? A custom type? What makes an
   empty cell distinct from an occupied one, and which choice makes an invalid symbol
   impossible to construct?

3. **Who owns "whose turn is it"?** The `Board`? The `Game`? A `Player`? Justify it —
   this is the question that separates a clean design from a tangled one.

4. **How does a move fail?** Go has no exceptions. What does your move method return,
   and how does a caller tell "cell taken" apart from "out of bounds"?

5. **Where does win-checking live?** Consider: if tomorrow the rules changed to
   "4-in-a-row on a 5x5 board wins," how many files would you have to touch? If the
   answer is more than one, redesign.

6. **What's the public API?** Write the method signatures a caller outside the package
   would use. That's your contract — everything else stays unexported.

---

## Follow-ups (the interviewer will ask these)

- Support an `N x N` board with a configurable win-length `K`.
- Add a third player. What breaks?
- Add an undo.
- Add a computer opponent with pluggable difficulty. Which pattern does this want?
- Make it safe for two goroutines submitting moves concurrently.

---

## Status

- [x] Design sketched
- [x] Types defined
- [x] Move validation
- [x] Win / draw detection
- [x] Tests passing (`go test ./problems/tictactoe/ -v`)
- [x] Follow-up: configurable `N` and `K`
- [x] Follow-up: three or more players
- [ ] Follow-up: undo
- [ ] Follow-up: computer opponent (Strategy)
- [ ] Follow-up: concurrent-safe moves

## How it fits together

| File | Responsibility |
|---|---|
| `symbol.go` | `Symbol` — a defined type so a stray character can't be a player's sign |
| `player.go` | `Player` — name and sign, no game state |
| `board.go` | Grid storage, bounds/occupancy validation, line detection |
| `game.go` | Turn order, game status, winner — the rules |
| `errors.go` | Sentinel errors so callers can tell *why* a move failed |

The split that matters: **`Board` never learns that players exist.** It takes a
`Symbol`, not a `Player`. That's what lets `Game` own turn order without
`Board` knowing there are two of anything — and it's why adding a third player
changed nothing in `board.go`.

`Move` runs its five steps in a fixed order, and the order *is* the design:

1. Reject if the game is already over
2. Let `Board` validate position and occupancy
3. Check for a win — **before** the draw check
4. Check for a draw
5. Advance the turn, only on a valid move that didn't end the game

Step 3 before step 4 is the subtle one: a final mark that fills the board *and*
completes a line is a win, not a draw. `TestWinningMoveThatAlsoFillsBoard`
pins that down.

## Try it

```go
x := tictactoe.NewPlayer("Xavier", tictactoe.X)
o := tictactoe.NewPlayer("Olivia", tictactoe.O)

g, err := tictactoe.NewGame(3, 3, x, o)   // 3x3, three in a row wins
if err != nil {
    log.Fatal(err)
}

if err := g.Move(0, 0); err != nil {
    fmt.Println("rejected:", err)
}

fmt.Print(g.Board())
if w, ok := g.Winner(); ok {
    fmt.Println(w, "wins")
}
```
