# lgo
http service frame golang cgi服务框架

## Example
```go
   package main
   
   import (
    "lgo"
   )
   
   func main() {
    lgo.Run(10168)
   }  
```

```go
  package main
   
   import (
    "fmt"
    "lgo"
   )
   
   func init() {
    lgo.HandleFunc("/getLineDataAPI", getLineInfo)
  }
  
  func getLineInfo(ctx *lgo.Context) {
   fmt.Println("adLineInfo")
  }
```
