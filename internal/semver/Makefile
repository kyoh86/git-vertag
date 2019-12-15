SUFFIX   = .y
SRCDIR   = ./

SOURCES  = $(wildcard $(SRCDIR)/*$(SUFFIX))
TARGETS  = $(notdir $(SOURCES:$(SUFFIX)=.y.go))


.SUFFIXES: .y .y.go
.PHONY: gen gen-parser gen-parser-force clean-parser

.y.y.go:
	@echo $(TARGETS)
	goyacc -p $(subst parse_,,$*) -o $@ $<
	rm y.output
	gofmt -w $@ 

gen: gen-parser

gen-parser: $(TARGETS)

gen-parser-force: clean-parser gen-parser

clean-parser:
	rm -f parse_*.y.go

