package tictactoe

// Symbol is a mark that can occupy a cell. Defining it as its own type (rather
// than using a bare rune) means the compiler rejects a stray character where a
// player's sign is expected.
type Symbol rune

const (
	// Empty is the zero-mark: a cell nobody has claimed.
	Empty Symbol = '-'
	X     Symbol = 'X'
	O     Symbol = 'O'
)

// String makes Symbol satisfy fmt.Stringer, so printing a cell shows "X"
// rather than its numeric code point.
func (s Symbol) String() string {
	return string(rune(s))
}
