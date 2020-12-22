# go-targetd

Go binding of [open-iscsi/targetd](https://github.com/open-iscsi/targetd) API.

## Usage

`targetd.yaml`

```yaml
user: "foo"
password: bar
ssl: false
```

```go
package main

import (
	"context"
	"fmt"
	"log"
	
	"github.com/lovi-cloud/go-targetd/targetd"
)

func main() {
	ctx := context.Background()
	
	client, err := targetd.New("http://192.0.2.1:18700", "foo", "bar", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	
	pools, err := client.GetPoolList(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", pools)
}
```