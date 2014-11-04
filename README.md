# go-namecheap

A Go library for using [the Namecheap API](https://www.namecheap.com/support/api/intro.aspx).

## Examples

```go
package main
import (
  "fmt"
  namecheap "github.com/billputer/go-namecheap"
)

func main() {
  apiUser := "billwiens"
  apiToken := "xxxxxxx"
  userName := "billwiens"

  client := namecheap.NewClient(apiUser, apiToken, userName)

  // Get a list of your domains
  domains, _ := client.Domains()
  for _, domain := range domains {
    fmt.Printf("Domain: %s\n", domain.Name)
  }

}
```

For more complete documentation, load up godoc and find the package.

## Development

- Source hosted at [GitHub](https://github.com/billputer/go-namecheap)
- Report issues and feature requests to [GitHub Issues](https://github.com/billputer/go-namecheap/issues)

Pull requests welcome!

## License

TBD