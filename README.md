# Inquestitobot

A component for search engines to create an inverted index. All data is
stored in an sqlite3 database (`inquestitobot.db`) using the defined
structure below.

## How to Use

`go build` && `./inquestitobot`

or

`./init`

## Stucture

Inquestitobot creates a `Document` which is a representation for a
resource (typically a webpage).

```Go
type Document struct {
	ID string
	Title string
	URL string
	Description string
	Checksum string
}
```

The `ID` variable is a checksum of the URL, `Title` is the title of the
resource, `URL` is the location of the resource, `Description` is a short
description of the resource, and `Checksum` is a checksum of the entire
`Document`.
