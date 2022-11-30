package main

import (
  "github.com/verdude/zapr"
)

func main() {
  zapr.Init(9)
  defer zapr.Sync()
  zapr.I("ok")
  zapr.V(8).I("um ok")
}
