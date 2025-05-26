# gobag
[![Go Reference](https://pkg.go.dev/badge/github.com/stianwa/gobag.svg)](https://pkg.go.dev/github.com/stianwa/gobag) [![Go Report Card](https://goreportcard.com/badge/github.com/stianwa/gobag)](https://goreportcard.com/report/github.com/stianwa/gobag)

Package gobag offers a collection of small, generic utility functions
for slices, maps, string parsing, and more.  It provides practical
helpers like deduplication, key extraction, and balanced string
splitting.

Installation
------------

The recommended way to install gobag

```
go get github.com/stianwa/gobag
```

Examples
--------

```go

package main
 
import (
        "fmt"
        "github.com/stianwa/gobag"
		"strings"
)

func main() {
        ds := gobag.Deduplicate([]string{"foo","bar","zot","bar"})                                  
        fmt.Printf("%s\n", strings.Join(ds,", "))
}
```

License
-------

MIT, see [LICENSE.md](LICENSE.md)
