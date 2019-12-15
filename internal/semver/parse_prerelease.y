%{
package semver

import "fmt"

%}

%union{
  pre PreRelease
  pre_ids []PreReleaseID
  pre_id PreReleaseID
  string []byte
  char byte
}

%type <pre> pre_release
%type <pre_ids> dot_separated_pre_release_identifiers
%type <pre_id> pre_release_identifier
%type <string> alphanum_identifier identifier_characters number digits
%type <char> digit non_digit identifier_character

%token <char> PRE_HYPHEN PRE_LETTER PRE_ZERO PRE_POSITIVE_DIGIT

%%

pre_release:
  dot_separated_pre_release_identifiers
    {
      prereleaselex.(*prereleaseLexerImpl).result = PreRelease($1)
    }

dot_separated_pre_release_identifiers:
  pre_release_identifier
    {
      $$ = []PreReleaseID{$1}
    }
| pre_release_identifier '.' dot_separated_pre_release_identifiers
    {
      $$ = append([]PreReleaseID{$1}, $3...)
    }

pre_release_identifier:
  alphanum_identifier
    {
      $$ = PreReleaseID{ str: string($1), isNum: false }
    }
| number
    {
      $$ = PreReleaseID{ str: string($1), num: numbytes($1), isNum: true }
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
  PRE_LETTER
| PRE_HYPHEN

number:
  PRE_ZERO
    { $$ = []byte{$1} }
| PRE_POSITIVE_DIGIT
    { $$ = []byte{$1} }
| PRE_POSITIVE_DIGIT number
    { $$ = append([]byte{$1}, $2...) }

digits:
  digit
    { $$ = []byte{$1} }
| digit digits
    { $$ = append([]byte{$1}, $2...) }

digit:
  PRE_ZERO
| PRE_POSITIVE_DIGIT

%%

type prereleaseLexerImpl struct {
  index int
  len int
  raw []rune
  syntaxErrPos int
  syntaxErr string
  result PreRelease
}

func (s *prereleaseLexerImpl) Lex(lval *prereleaseSymType) int {
  if s.len <= s.index {
    return 0
  }
  n := s.raw[s.index]
  s.index++
  switch n {
  /* PRE_HYPHEN: '-' */
  case '-':
    lval.char = byte(n)
    return PRE_HYPHEN

  /* PRE_ZERO: '0' */
  case '0':
    lval.char = byte(n)
    return PRE_ZERO
  }

  /* PRE_POSITIVE_DIGIT: '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' */
  if '1' <= n && n <= '9' {
    lval.char = byte(n)
    return PRE_POSITIVE_DIGIT
  }

  /* PRE_LETTER:
    'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J'
  | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T'
  | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | 'a' | 'b' | 'c' | 'd'
  | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n'
  | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x'
  | 'y' | 'z' */
  if 'A' <= n && n <= 'Z' || 'a' <= n && n <= 'z' {
    lval.char = byte(n)
    return PRE_LETTER
  }

  // TODO; check valid char
  // other
  return int(n)
}

func (s *prereleaseLexerImpl) Rune(r rune) int {
  return -1
}

func (s *prereleaseLexerImpl) Error(err string) {
  s.syntaxErrPos = s.index-1
  s.syntaxErr = err
}

func MustParsePreRelease(s string) PreRelease {
  pre, err := ParsePreRelease(s)
  if err != nil {
    panic(err)
  }
  return pre
}

func ParsePreRelease(s string) (PreRelease, error) {
  lex := &prereleaseLexerImpl{ raw: []rune(s), len: len(s) }
  prereleaseParse(lex)
  if lex.syntaxErr != "" {
    return lex.result, fmt.Errorf("prerelease %s at %d", lex.syntaxErr, lex.syntaxErrPos)
  }
  return lex.result, nil
}

