[![progress-banner](https://backend.codecrafters.io/progress/grep/ba3e1c9e-bfcb-4eb7-9891-96aaefaa5967)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

This is a starting point for Go solutions to the
["Build Your Own grep" Challenge](https://app.codecrafters.io/courses/grep/overview).

[Regular expressions](https://en.wikipedia.org/wiki/Regular_expression)
(Regexes, for short) are patterns used to match character combinations in
strings. [`grep`](https://en.wikipedia.org/wiki/Grep) is a CLI tool for
searching using Regexes.

In this challenge you'll build your own implementation of `grep`. Along the way
we'll learn about Regex syntax, how parsers/lexers work, and how regular
expressions are evaluated.

**Note**: If you're viewing this repo on GitHub, head over to
[codecrafters.io](https://codecrafters.io) to try the challenge.

# Passing the first stage

The entry point for your `grep` implementation is in `app/main.go`. Study and
uncomment the relevant code, and push your changes to pass the first stage:

```sh
git commit -am "pass 1st stage" # any msg
git push origin master
```

Time to move on to the next stage!

# Features Implemented

This grep implementation includes:
- Basic character matching and character classes (\d, \w, [abc], [^xyz])
- Anchors (^, $)
- Quantifiers (+, ?, *)
- Wildcard (.)
- Alternation (|)
- Backreferences (\1-\9) with nested group support
- File search (single, multiple, multi-line)
- Recursive directory search (-r flag)

# Stage 2 & beyond

Note: This section is for stages 2 and beyond.

1. Ensure you have `go (1.24)` installed locally
1. Run `./your_program.sh` to run your program, which is implemented in
   `app/main.go`.
1. Commit your changes and run `git push origin master` to submit your solution
   to CodeCrafters. Test output will be streamed to your terminal.
