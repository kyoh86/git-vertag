%{
package semver

import (
  "fmt"
  "errors"
)

%}

%union{
  ver Semver
  pre PreRelease
  pre_ids []PreReleaseID
  pre_id PreReleaseID
  bld Build
  bds []BuildID
  bid BuildID
  string []byte
  char byte
}

%type <ver> valid_semver
%type <pre> pre_release
%type <pre_ids> dot_separated_pre_release_identifiers
%type <pre_id> pre_release_identifier
%type <bld> build
%type <bds> dot_separated_build_identifiers
%type <bid> build_identifier

%type <string> alphanum_identifier
%type <string> identifier_characters

%type <string> number major minor patch
%type <string> digits
%type <char> identifier_character digit non_digit

%token <char> SEMVER_HYPHEN SEMVER_LETTER SEMVER_ZERO SEMVER_POSITIVE_DIGIT

%%

valid_semver:
  major '.' minor '.' patch
    {
      semverlex.(*semverLexerImpl).result = Semver{ Major: numbytes($1), Minor: numbytes($3), Patch: numbytes($5) }
    }
| major '.' minor '.' patch SEMVER_HYPHEN pre_release
    {
      semverlex.(*semverLexerImpl).result = Semver{ Major: numbytes($1), Minor: numbytes($3), Patch: numbytes($5), PreRelease: $7 }
    }
| major '.' minor '.' patch '+' build
    {
      semverlex.(*semverLexerImpl).result = Semver{ Major: numbytes($1), Minor: numbytes($3), Patch: numbytes($5), Build: $7 }
    }
| major '.' minor '.' patch SEMVER_HYPHEN pre_release '+' build
    {
      semverlex.(*semverLexerImpl).result = Semver{ Major: numbytes($1), Minor: numbytes($3), Patch: numbytes($5), PreRelease: $7, Build: $9 }
    }

major:
  number
minor:
  number
patch:
  number

pre_release:
  dot_separated_pre_release_identifiers
    {
      $$ = PreRelease($1)
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

build:
  dot_separated_build_identifiers
     {
        $$ = Build($1)
     }
dot_separated_build_identifiers:
  build_identifier
    {
      $$ = []BuildID{$1}
    }
| build_identifier '.' dot_separated_build_identifiers
    {
      $$ = append([]BuildID{$1}, $3...)
    }

build_identifier:
  alphanum_identifier
    {
      $$ = BuildID(string($1))
    }
| digits
    {
      $$ = BuildID(string($1))
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
  SEMVER_LETTER
| SEMVER_HYPHEN

number:
  SEMVER_ZERO
    { $$ = []byte{$1} }
| SEMVER_POSITIVE_DIGIT
    { $$ = []byte{$1} }
| SEMVER_POSITIVE_DIGIT number
    { $$ = append([]byte{$1}, $2...) }

digits:
  digit
    { $$ = []byte{$1} }
| digit digits
    { $$ = append([]byte{$1}, $2...) }

digit:
  SEMVER_ZERO
| SEMVER_POSITIVE_DIGIT

%%

type semverLexerImpl struct {
  pos int
  len int

  bytes []byte
  err string

  result Semver
}

var lexm = map[lex]int {
  DOT: int('.'),
  PLUS: int('+'),
	HYPHEN: SEMVER_HYPHEN,
	LETTER: SEMVER_LETTER,
	POSITIVE: SEMVER_POSITIVE_DIGIT,
	ZERO  : SEMVER_ZERO,
}

func (s *semverLexerImpl) Lex(lval *semverSymType) int {
  if s.len <= s.pos {
    return 0
  }
  n := s.bytes[s.pos]
  s.pos++
  lval.char = n

	l, ok := lexMap[n]
	if !ok {
    s.err = fmt.Sprintf("invalid char at %d", s.pos-1)
		return int(n)
	}
  return lexm[l]
}

func (s *semverLexerImpl) Error(err string) {
  s.err = err
}

func Parse(s string) (Semver, error) {
  lex := &semverLexerImpl{ bytes: []byte(s), len: len(s) }
  semverParse(lex)
  if lex.err != "" {
    return lex.result, errors.New(lex.err)
  }
  return lex.result, nil
}
