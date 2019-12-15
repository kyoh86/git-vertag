%{
package semver

import "fmt"

%}

%union{
  pri PreReleaseID
  string []byte
  char byte
}

%type <pri> pri_release_identifier
%type <string> alphanum_identifier identifier_characters number digits
%type <char> digit non_digit identifier_character

%token <char> PRI_HYPHEN PRI_LETTER PRI_ZERO PRI_POSITIVE_DIGIT

%%

pri_release_identifier:
  alphanum_identifier
    {
      prilex.(*priLexerImpl).result = PreReleaseID{ str: string($1), isNum: false }
    }
| number
    {
      prilex.(*priLexerImpl).result = PreReleaseID{ str: string($1), num: numbytes($1), isNum: true }
    }

alphanum_identifier:
  non_digit
    { $$ = []byte{$1} }
| non_digit identifier_characters
    { $$ = append([]byte{$1}, $2...) }
| identifier_characters non_digit
    { $$ = append($1, $2) }
| identifier_characters non_digit identifier_characters
    { $$ = append(append($1, $2), $3...) }

identifier_characters:
  identifier_character
    { $$ = []byte{$1} }
| identifier_character identifier_characters
    { $$ = append([]byte{$1}, $2...) }

identifier_character:
  digit
| non_digit

non_digit:
  PRI_LETTER
| PRI_HYPHEN

number:
  PRI_ZERO
    { $$ = []byte{$1} }
| PRI_POSITIVE_DIGIT
    { $$ = []byte{$1} }
| PRI_POSITIVE_DIGIT number
    { $$ = append([]byte{$1}, $2...) }

digits:
  digit
    { $$ = []byte{$1} }
| digit digits
    { $$ = append([]byte{$1}, $2...) }

digit:
  PRI_ZERO
| PRI_POSITIVE_DIGIT

%%

type priLexerImpl struct {
  index int
  len int
  raw []rune
  syntaxErrPos int
  syntaxErr string
  result PreReleaseID
}

func (s *priLexerImpl) Lex(lval *priSymType) int {
  if s.len <= s.index {
    return 0
  }
  n := s.raw[s.index]
  s.index++
  switch n {
  /* PRI_HYPHEN: '-' */
  case '-':
    lval.char = byte(n)
    return PRI_HYPHEN

  /* PRI_ZERO: '0' */
  case '0':
    lval.char = byte(n)
    return PRI_ZERO
  }

  /* PRI_POSITIVE_DIGIT: '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' */
  if '1' <= n && n <= '9' {
    lval.char = byte(n)
    return PRI_POSITIVE_DIGIT
  }

  /* PRI_LETTER:
    'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J'
  | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T'
  | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | 'a' | 'b' | 'c' | 'd'
  | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n'
  | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x'
  | 'y' | 'z' */
  if 'A' <= n && n <= 'Z' || 'a' <= n && n <= 'z' {
    lval.char = byte(n)
    return PRI_LETTER
  }

  // TODO; check valid char
  // other
  return int(n)
}

func (s *priLexerImpl) Rune(r rune) int {
  return -1
}

func (s *priLexerImpl) Error(err string) {
  s.syntaxErrPos = s.index-1
  s.syntaxErr = err
}

func MustParsePreReleaseID(s string) PreReleaseID {
  pri, err := ParsePreReleaseID(s)
  if err != nil {
    panic(err)
  }
  return pri
}

func ParsePreReleaseID(s string) (PreReleaseID, error) {
  lex := &priLexerImpl{ raw: []rune(s), len: len(s) }
  priParse(lex)
  if lex.syntaxErr != "" {
    return lex.result, fmt.Errorf("%s at %d", lex.syntaxErr, lex.syntaxErrPos)
  }
  return lex.result, nil
}

