package semver

const (
	B_ZERO   = byte('0')
	B_ONE    = byte('1')
	B_NINE   = byte('9')
	B_UP_A   = byte('A')
	B_UP_Z   = byte('Z')
	B_LOW_A  = byte('a')
	B_LOW_Z  = byte('z')
	B_HYPHEN = byte('-')
	B_DOT    = byte('.')
	B_PLUS   = byte('+')
)

func numbytes(b []byte) (i uint64) {
	for _, x := range b {
		i = i*10 + uint64(x-B_ZERO)
	}
	return
}

const (
	/* HYPHEN: '-' */
	HYPHEN int = iota
	LETTER
	POSITIVE
	ZERO
	DOT  int = int('.')
	PLUS int = int('+')
)

var lexMap = map[byte]int{
	byte('-'): HYPHEN,
	byte('0'): ZERO,
	byte('1'): POSITIVE,
	byte('2'): POSITIVE,
	byte('3'): POSITIVE,
	byte('4'): POSITIVE,
	byte('5'): POSITIVE,
	byte('6'): POSITIVE,
	byte('7'): POSITIVE,
	byte('8'): POSITIVE,
	byte('9'): POSITIVE,
	byte('A'): LETTER,
	byte('B'): LETTER,
	byte('C'): LETTER,
	byte('D'): LETTER,
	byte('E'): LETTER,
	byte('F'): LETTER,
	byte('G'): LETTER,
	byte('H'): LETTER,
	byte('I'): LETTER,
	byte('J'): LETTER,
	byte('K'): LETTER,
	byte('L'): LETTER,
	byte('M'): LETTER,
	byte('N'): LETTER,
	byte('O'): LETTER,
	byte('P'): LETTER,
	byte('Q'): LETTER,
	byte('R'): LETTER,
	byte('S'): LETTER,
	byte('T'): LETTER,
	byte('U'): LETTER,
	byte('V'): LETTER,
	byte('W'): LETTER,
	byte('X'): LETTER,
	byte('Y'): LETTER,
	byte('Z'): LETTER,
	byte('a'): LETTER,
	byte('b'): LETTER,
	byte('c'): LETTER,
	byte('d'): LETTER,
	byte('e'): LETTER,
	byte('f'): LETTER,
	byte('g'): LETTER,
	byte('h'): LETTER,
	byte('i'): LETTER,
	byte('j'): LETTER,
	byte('k'): LETTER,
	byte('l'): LETTER,
	byte('m'): LETTER,
	byte('n'): LETTER,
	byte('o'): LETTER,
	byte('p'): LETTER,
	byte('q'): LETTER,
	byte('r'): LETTER,
	byte('s'): LETTER,
	byte('t'): LETTER,
	byte('u'): LETTER,
	byte('v'): LETTER,
	byte('w'): LETTER,
	byte('x'): LETTER,
	byte('y'): LETTER,
	byte('z'): LETTER,
	byte('.'): DOT,
	byte('+'): PLUS,
}

type lexerBase struct {
	pos int
	len int

	lex   int
	ok    bool
	cur   byte
	bytes []byte
	err   string

	lexm map[int]int
}

func (s *lexerBase) chLex() (byte, int) {
	if s.len <= s.pos {
		return 0, 0
	}
	s.cur = s.bytes[s.pos]
	s.pos++
	s.lex, s.ok = lexMap[s.cur]
	if !s.ok {
		return s.cur, int(s.cur)
	}
	return s.cur, s.lexm[s.lex]
}

func (s *lexerBase) Error(err string) {
	s.err = err
}
