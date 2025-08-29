# SQL Parser Setup Guide

This guide explains how to set up the SQL parser for Go using ANTLR4.

## Prerequisites

- **Install ANTLR4 tools**  
    ```sh
    pip install antlr4-tools
    ```

## Download Grammar Files

- Download your parser and lexer from the [ANTLR grammars-v4 repository](https://github.com/antlr/grammars-v4/tree/master).  
    **Note:** Please review the license of these grammar files before use.

## Generate Go Parser Code

- Run the following command to generate Go code with visitor support:
    ```sh
    antlr4 -Dlanguage=Go -visitor -package grammar -o grammar YOURLexer.g4 YOURParser.g4
    ```
    - `-Dlanguage=Go`: Generates Go code.
    - `-visitor`: Enables visitor pattern generation.
    - `-package grammar`: Sets the Go package name to `grammar`.
    - `-o grammar`: Outputs generated files to the `grammar` directory.

## References

- [ANTLR grammars-v4 repository](https://github.com/antlr/grammars-v4/tree/master)
- [ANTLR4 Go Target Documentation](https://github.com/antlr/antlr4/blob/master/doc/go-target.md)
