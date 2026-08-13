# Low Level Design in Go

Learning Low Level Design (LLD) / Object-Oriented Design by building every concept
from scratch in Go — a language with **no classes and no inheritance**.

That constraint is the point. Go forces you to express object-oriented design with
**interfaces, struct embedding, and composition**, which means you learn *why* a
pattern exists instead of memorizing its Java boilerplate.

> This is a learning journal. Every file here was written by hand, one topic at a time.

---

## Roadmap

### 1. OOP foundations in Go
- [ ] Structs, methods, and receivers (value vs pointer)
- [ ] Interfaces and implicit satisfaction
- [ ] Encapsulation via package boundaries and exported/unexported identifiers
- [ ] Abstraction with interfaces
- [ ] Composition and struct embedding (Go's answer to inheritance)
- [ ] Polymorphism via interfaces
- [ ] Enums with `iota` and typed constants

### 2. Class relationships
- [ ] Association
- [ ] Aggregation
- [ ] Composition
- [ ] Dependency

### 3. Design principles
- [ ] DRY, KISS, YAGNI
- [ ] **S**ingle Responsibility
- [ ] **O**pen/Closed
- [ ] **L**iskov Substitution
- [ ] **I**nterface Segregation
- [ ] **D**ependency Inversion

### 4. Design patterns

**Creational**
- [ ] Singleton · [ ] Factory Method · [ ] Abstract Factory · [ ] Builder · [ ] Prototype

**Structural**
- [ ] Adapter · [ ] Bridge · [ ] Composite · [ ] Decorator · [ ] Facade · [ ] Flyweight · [ ] Proxy

**Behavioral**
- [ ] Strategy · [ ] Observer · [ ] Command · [ ] State · [ ] Iterator · [ ] Template Method
- [ ] Chain of Responsibility · [ ] Mediator · [ ] Memento · [ ] Visitor

### 5. Concurrency
- [ ] Goroutines and the scheduler
- [ ] Channels, `select`, and directional channels
- [ ] `sync.Mutex`, `RWMutex`, `WaitGroup`, `Once`
- [ ] Race conditions and `go test -race`
- [ ] Worker pools and fan-in / fan-out
- [ ] Deadlock and how Go detects it

### 6. Interview problems

- [ ] [Tic Tac Toe](problems/tictactoe)
- [ ] Parking Lot
- [ ] Vending Machine
- [ ] Logging Framework
- [ ] Stack Overflow
- [ ] ATM
- [ ] LRU Cache
- [ ] Elevator System
- [ ] Splitwise
- [ ] Hotel Management
- [ ] Rate Limiter
- [ ] Ride Sharing (Uber)
- [ ] Music Streaming (Spotify)
- [ ] Online Shopping (Amazon)
- [ ] Food Delivery

---

## Layout

```
oop/          fundamentals — structs, interfaces, embedding, polymorphism
principles/   SOLID, DRY, KISS, YAGNI — each with a bad/ and good/ example
patterns/     creational, structural, behavioral
concurrency/  goroutines, channels, sync primitives
problems/     one folder per interview problem — self-contained package + README
```

Every directory is a Go package with a `README.md` explaining the concept, the code
itself, and a `_test.go` proving it works.

## Running

```bash
go test ./...
```

## Progress

Started: 2026-08-13
