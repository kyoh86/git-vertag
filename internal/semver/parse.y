%{
package semver

import "fmt"

type num struct {
  num int
  str string
}
%}

%union{
  ver Semver
  pre PreRelease
  pds []PreReleaseID
  pid PreReleaseID
  bld Build
  bds []BuildID
  bid BuildID
  str string
  num num
}

%type <ver> valid_semver
%type <pre> pre_release
%type <pds> dot_separated_pre_release_identifiers
%type <pid> pre_release_identifier
%type <bld> build
%type <bds> dot_separated_build_identifiers
%type <bid> build_identifier

%type <str> alphanum_identifier
%type <str> identifier_character identifier_characters non_digit

%type <num> major minor patch number
%type <num> digit digits

%token <str> HYPHEN LETTER
%token <num> ZERO POSITIVE_DIGIT

%%

valid_semver:
  major '.' minor '.' patch
    {
      yylex.(*semverLexer).result = Semver{ Major: uint64($1.num), Minor: uint64($3.num), Patch: uint64($5.num) }
    }
| major '.' minor '.' patch HYPHEN pre_release
    {
      yylex.(*semverLexer).result = Semver{ Major: uint64($1.num), Minor: uint64($3.num), Patch: uint64($5.num), PreRelease: $7 }
    }
| major '.' minor '.' patch '+' build
    {
      yylex.(*semverLexer).result = Semver{ Major: uint64($1.num), Minor: uint64($3.num), Patch: uint64($5.num), Build: $7 }
    }
| major '.' minor '.' patch HYPHEN pre_release '+' build
    {
      yylex.(*semverLexer).result = Semver{ Major: uint64($1.num), Minor: uint64($3.num), Patch: uint64($5.num), PreRelease: $7, Build: $9 }
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

pre_release_identifier:
  alphanum_identifier
    {
      $$ = PreReleaseID{ str: $1, isNum: false }
    }
| number
    {
      $$ = PreReleaseID{ str: $1.str, num: uint64($1.num), isNum: true }
    }

build_identifier:
  alphanum_identifier
    {
      $$ = BuildID($1)
    }
| digits
    {
      $$ = BuildID($1.str)
    }

alphanum_identifier:
  non_digit
| non_digit identifier_characters
    { $$ = $1+$2 }
| identifier_characters non_digit
    { $$ = $1+$2 }
| identifier_characters non_digit identifier_characters
    { $$ = $1+$2+$3 }

identifier_characters:
  identifier_character
| identifier_character identifier_characters
    { $$ = $1+$2 }

identifier_character:
  digit
    { $$ = $1.str }
| non_digit

non_digit:
  LETTER
| HYPHEN

number:
  ZERO
| POSITIVE_DIGIT
| POSITIVE_DIGIT digits
    { $$ = num{ num: $1.num*10 + $2.num, str: $1.str + $2.str } }

digits:
  digit
| digit digits
    { $$ = num{ num: $1.num*10 + $2.num, str: $1.str + $2.str } }

digit:
  ZERO
| POSITIVE_DIGIT

%%

type semverLexer struct {
  index int
  len int
  raw []rune
  syntaxErrPos int
  syntaxErr string
  result Semver
}

func (s *semverLexer) Lex(lval *yySymType) int {
  if s.len <= s.index {
    return 0
  }
  n := s.raw[s.index]
  s.index++
  switch n {
    /* HYPHEN: '-' */
    case '-':
      return HYPHEN
    /* ZERO: '0' */
    case '0':
      lval.num = num{ num: 0, str: "0" }
      return ZERO
    /* POSITIVE_DIGIT: '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' */
    case '1', '2', '3', '4', '5', '6', '7', '8', '9':
      lval.num = num{ num: int(n-'0'), str: string([]rune{n}) }
      return POSITIVE_DIGIT
    /* LETTER:
      'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J'
    | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T'
    | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | 'a' | 'b' | 'c' | 'd'
    | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n'
    | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x'
    | 'y' | 'z' */
    case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J':
      lval.str = string([]rune{n})
      return LETTER
    case 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T':
      lval.str = string([]rune{n})
      return LETTER
    case 'U', 'V', 'W', 'X', 'Y', 'Z', 'a', 'b', 'c', 'd':
      lval.str = string([]rune{n})
      return LETTER
    case 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n':
      lval.str = string([]rune{n})
      return LETTER
    case 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x':
      lval.str = string([]rune{n})
      return LETTER
    case 'y', 'z':
      lval.str = string([]rune{n})
      return LETTER
    default:
      return int(n)
  }
}

func (s *semverLexer) Rune(r rune) int {
  return -1
}

func (s *semverLexer) Error(err string) {
  s.syntaxErrPos = s.index-1
  s.syntaxErr = err
}

func Parse(s string) (Semver, error) {
  lex := &semverLexer{ raw: []rune(s), len: len(s) }
  yyParse(lex)
  if lex.syntaxErr != "" {
    return lex.result, fmt.Errorf("%s at %d", lex.syntaxErr, lex.syntaxErrPos)
  }
  return lex.result, nil
}
