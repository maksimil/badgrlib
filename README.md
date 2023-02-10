# badrlib

## Spec

### Types

- `Format`

  Contains info on how to to place info on a badge and
  how many badges to fit on a page (the type declaration is in `format.go`)

  `Format` can also be parsed from a toml file (`example/format.toml`)

- `InputTable`

  Contains info to put on the badges (`format.go`)

  `InputTable` can also be parsed from a csv file (`example/names.csv`)

- `Context`

  Contains canvas context and a font family (from `tdewolff/canvas`)

- `ContextDrawer: func (Context)`

  Is a function that draws on a `Context`

### Functions

- `ParseFormat: string -> (Format, error)`

  Returns an error on toml parsing errors

- `ParseTable: string -> (InputTable, error)`

  Input string can be either LF or CRLF

  Returns an error if one line contains a number of elements different
  from the first line (empty lines however are allowed)

- `CreateObjectDrawer: (Format, map[string]string) -> ContextDrawer`

  Returns a `ContextDrawer` that draws

- `FitObjectsOnPage: (Format, []ContextDrawer) -> ContextDrawer`

  Returns a complete svg containing the first elements from the list to fit
  on the paper
