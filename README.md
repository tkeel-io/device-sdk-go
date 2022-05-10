# tkeel device sdk for go

## Description

### Organization --TODO

- __mqtt__ - [paho.mqtt]()
- __examples__ - some example use this sdk
- __tkeel__ - ues mqtt to 
- __util__ - some useful util

### Quick start -- TODO:


```go
package main

import (
    cli "github.com/tkeel-io/device-sdk-go"
)

func main() {
    client, err := cli.NewClient()
    if err != nil {
        panic(err)
    }
    defer client.Close()
    //TODO: use the client here, see below for examples 
}
```



## Usage
> Assuming you already have [installed](https://golang.org/doc/install) Go

