%{
package semver

import "fmt"

%}

%union{
  bid BuildID
  string []byte
  char byte
}

%type <bid> build_identifier

%type <string> alphanum_identifier
%type <string> identifier_characters

%type <string> number
%type <string> digits
%type <char> identifier_character digit non_digit

%token <char> BUILD_HYPHEN BUILD_LETTER BUILD_ZERO BUILD_POSITIVE_DIGIT

%%

build_identifier:
  alphanum_identifier
    {
      bdilex.(*bdiLexerImpl).result = BuildID(string($1))
    }
| digits
    {
      bdilex.(*bdiLexerImpl).result = BuildID(string($1))
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
  BUILD_LETTER
| BUILD_HYPHEN

number:
  BUILD_ZERO
    { $$ = []byte{$1} }
| BUILD_POSITIVE_DIGIT
    { $$ = []byte{$1} }
| BUILD_POSITIVE_DIGIT number
    { $$ = append([]byte{$1}, $2...) }

digits:
  digit
    { $$ = []byte{$1} }
| digit digits
    { $$ = append([]byte{$1}, $2...) }

digit:
  BUILD_ZERO
| BUILD_POSITIVE_DIGIT

%%

type bdiLexerImpl struct {
  index int
  len int
  raw []rune
  err error
  syntaxErrPos int
  syntaxErr string
  result BuildID
}

func (s *bdiLexerImpl) Lex(lval *bdiSymType) int {
  if s.err != nil {
    return 0
  }
  if s.len <= s.index {
    return 0
  }
  n := s.raw[s.index]
  s.index++
  lval.char = byte(n)
  switch n {
  /* BUILD_HYPHEN: '-' */
  case '-':
    return BUILD_HYPHEN

  /* BUILD_ZERO: '0' */
  case '0':
    return BUILD_ZERO
  }

  /* BUILD_POSITIVE_DIGIT: '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' */
  if '1' <= n && n <= '9' {
    return BUILD_POSITIVE_DIGIT
  }

  /* BUILD_LETTER:
    'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J'
  | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T'
  | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | 'a' | 'b' | 'c' | 'd'
  | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n'
  | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x'
  | 'y' | 'z' */
  if 'A' <= n && n <= 'Z' || 'a' <= n && n <= 'z' {
    return BUILD_LETTER
  }

  // TODO; check valid char
  if n == '.'  || n == '+' {
    return int(n)
  }
  // other
  s.err = fmt.Errorf("invalid char at %d", s.index-1)
  return int(n)
}

func (s *bdiLexerImpl) Rune(r rune) int {
  return -1
}

func (s *bdiLexerImpl) Error(err string) {
  s.syntaxErrPos = s.index-1
  s.syntaxErr = err
}

func MustParseBuildID(s string) BuildID {
  pre, err := ParseBuildID(s)
  if err != nil {
    panic(err)
  }
  return pre
}

func ParseBuildID(s string) (BuildID, error) {
  lex := &bdiLexerImpl{ raw: []rune(s), len: len(s) }
  bdiParse(lex)
  if lex.err != nil {
    return lex.result, lex.err
  }
  if lex.syntaxErr != "" {
    return lex.result, fmt.Errorf("%s at %d", lex.syntaxErr, lex.syntaxErrPos)
  }
  return lex.result, nil
}
