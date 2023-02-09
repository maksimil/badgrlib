# badrlib

## Spec

### Format

`Format` type includes info on how to to place info on a badge and
how many badges to fit on a page (the type declaration is in `format.go`)

`Format` can also be parsed from a toml file (`example/format.toml`)

`InputTable` type includes info to put on the badges (`format.go`)

`InputTable` can also be parsed from a csv file (`example/names.csv`)

### Functions

- `ParseFormat: string -> (Format, error)`

  Returns an error on toml parsing errors

- `ParseTable: string -> (InputTable, error)`

  Input string can be either LF or CRLF

  Returns an error if one line contains a number of elements different
  from the first line (empty lines however are allowed)

- `CreateSingleSvg: (Format, map[string]string) -> (string, error)`

  Does not return a valid svg, just a concatenation of text elements

  Returns an error on template errors

- `WrapSvg: (string, Dimensions) -> (string, error)`

  Wraps svg into an `<svg></svg>` element

  Returns an error on template errors

- `FitObjectsOnPage: (Format, []string) -> (string, error)`

  Returns a complete svg containing the first elements from the list to fit
  on the paper
