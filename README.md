# My go router

This is my implementation of router. I used trie data structure with hash table

#### How to install
You can download this module and use it in your project by executing next command

```shell
export GO111MODULE=on
go get github.com/AndriyAntonenko/goRouter
```

#### How to use
`goRouter` module working with standard `net/http` module

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    
    "github.com/AndriyAntonenko/goRouter"
)

func main() {
    router := goRouter.NewRouter()
    
    router.Get("/api/users/v1/:id/:token/check/:hash", func(w http.ResponseWriter, r *http.Request, ps *goRouter.RouterParams) {
        // Get path parameters
        fmt.Println(ps.GetString("token"))
        fmt.Println(ps.ParseInt("id"))
        fmt.Println(ps.GetString("hash"))
        
        fmt.Fprintf(w, "Something")
    })

    server := http.Server{
        Addr:           ":8080",
        Handler:        router,
        MaxHeaderBytes: 1 << 20,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
	  }
    
    // ... Run the server
}
```

#### Future
In the next version I'm going to implement `catchAll` routes
