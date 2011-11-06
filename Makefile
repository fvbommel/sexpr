# This file is subject to a 1-clause BSD license.
# Its contents can be found in the enclosed LICENSE file.

include $(GOROOT)/src/Make.inc

TARG = github.com/jteeuwen/sexpr
GOFILES = error.go syntax.go token.go lexer.go ast.go node.go parse.go

include $(GOROOT)/src/Make.pkg
