
# parser/tokenizer.go assigns yySymType fields directly (lval.token = ...), so goyacc must emit a
# struct with typed fields. The upstream x/tools build does; some third-party builds (e.g. one
# installed as ~/bin/goyacc) emit `union any` + accessor methods instead and break the build.
# Prefer ~/go/bin/goyacc explicitly: PATH order differs per shell (fish finds ~/bin first,
# bash ~/go/bin first), so relying on plain `goyacc` makes the recipe shell-dependent.
# Override with GOYACC=/path/to/goyacc if needed.
goyacc_default := home_directory() / "go/bin/goyacc"
goyacc := env("GOYACC", if path_exists(goyacc_default) == "true" { goyacc_default } else { "goyacc" })

set shell := ["bash", "-euo", "pipefail", "-c"]

sql:
    {{goyacc}} -l -o parser/sql.go -v parser/sql.y.output parser/sql.y

join:
   {{goyacc}} -l -o join.go -v join.y.output join.y

run:
    /usr/local/bin/flex -G -i -L -o parser/tokenizer.go -S parser/flex_go.skl parser/tokenizer.l
    sed 's/\byyDef\b/yyDefLexer/g' parser/tokenizer.go > parser/tokenizer.go1
    sed 's/\byyChk\b/yyChkLexer/g' parser/tokenizer.go1 > parser/tokenizer.go2
    sed 's/\/\/line.*//g' parser/tokenizer.go2 > parser/tokenizer.go
    rm parser/tokenizer.go1
    rm parser/tokenizer.go2
    {{goyacc}} -l -o parser/sql.go -v parser/sql.y.output parser/sql.y 
    go build data/*.go
    go build parser/*.go
    go build machine/*.go
    go build driver/*.go
    go build tests/*.go
    go build example.go
    ./example

