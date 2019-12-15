package semver

const (
	B_ZERO = byte('0')
)

func numbytes(b []byte) (i uint64) {
	for _, x := range b {
		i = i*10 + uint64(x-B_ZERO)
	}
	return
}

type lex int

const (
	/* HYPHEN: '-' */
	HYPHEN lex = iota
	LETTER
	POSITIVE
	ZERO
	DOT  lex = lex('.')
	PLUS lex = lex('+')
)

var lexMap = map[byte]lex{
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

	bytes []byte
	err   string
	lexm  map[lex]int
}

func (s *lexerBase) chLex() (byte, int) {
	if s.len <= s.pos {
		return 0, 0
	}
	n := s.bytes[s.pos]
	s.pos++
	l, ok := lexMap[n]
	if !ok {
		return n, int(n)
	}
	return n, s.lexm[l]
}

func (s *lexerBase) Error(err string) {
	s.err = err
}
