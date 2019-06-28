# ivent

Logs free debug library

- [ivent](#ivent)
  - [Motivation](#Motivation)
  - [Installation](#Installation)
  - [Usage](#Usage)
    - [For libraries using ivent Logger](#For-libraries-using-ivent-Logger)
    - [For users](#For-users)

## Motivation

Diferent libraries choose a logging library for debug and traceability. Ivent pretends to move logging from author to users

## Installation

```bash
    go get github.com/licaonfee/ivent
```

## Usage

### For libraries using ivent Logger

```go
package lib

import (
    "github.com/licaonfee/ivent/log"
    "github.com/licaonfee/ivent/stream"
    )

var logger *ivent.Logger

func init(){
    //By default all clients use stream.Noop
    logger = log.NewLogger()
}

func DoSomethingWithIvent(){
    //do the thing
    logger.Trace("I do the thing")
}
```

### For users

```go
package main

import (
    "log"
    "github.com/licaonfee/ivent/stream"
    "github.com/licaonfee/ivent/log"
    "thelibrary/lib"
    )

func main(){
    stream := stream.NewAsyncStream()
    lib.logger.WithStream(stream)

    go func(){
        for e:= <-stream.Get(){
            log.Printf("%s : %s",e.Class, e.Data)
        }
    }

    //Prints:"Trace : I do the thing"
    lib.DoSomethingWithIvent()
}
```
