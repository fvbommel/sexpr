##  sexpr (A Go S-expression parser)

This is a full rewrite and redesign of github.com/jteeuwen/gsx.
It is faster and has a considerably lower memory footprint. GSX will no
longer be developed and should not be used. Use this package instead.
GSX remains for legacy and compatibility reasons.

This package offers a configurable S-Expression parser. It takes any
input and turns it into a parse tree.

The lexer and AST builder retain filename/line/col information from the
original input source. This should help with debugging potential errors during
parsing or at any later stage.

Some parts can be customized before parsing, in order to alter the behaviour and
output of the lexer and AST builder. For this purpose you have to supply an
instance of `sexpr.Syntax` to the parse function. It has the following fields:

    // A set of list delimiters. These are pairs of strings denoting the
    // start and end of an S-expression.
    // E.g.: "(", ")"
    Delimiters [][2]string
    
    // This string starts a single line comment.
    // A single line comment runs until the end of a line.
    // E.g: "//"
    SingleLineComment string
    
    // These strings denote what a multi-line comment starts with
    // and ends with.
    // E.g.: "/*", "*/"
    MultiLineComment []string
    
    // These strings determine how a string literal starts and ends.
    // E.g.: "abc".
    StringLit []string
    
    // These strings determine how a raw string literal starts and ends.
    // A raw string does not have its escape sequences parsed.
    // E.g.: `abc`.
    RawStringLit []string
    
    // These strings determine how a char literal starts and ends.
    // E.g.: 'a'.
    CharLit []string
    
    // This function should return whether or not the given
    // input qualifies as a boolean.
    BooleanFunc SyntaxFunc
    
    // This function should return whether or not the given
    // input qualifies as a number.
    NumberFunc SyntaxFunc

A single AST tree can be used in multiple `Parse()` calls for different
source files. Their output will then be merged with the given AST.

For an example of how to use this package, refer to `sexpr_test.go`.

## Dependencies

None.

## Build

    $ goinstall github.com/jteeuwen/sexpr

    import "github.com/jteeuwen/sexpr"

## License

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.

