package main

import (
	_ "go.uber.org/automaxprocs"
	"math/rand"
	"minifast/app/user/srv"
	"time"
)

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	srv.NewApp("user-server").Run()
}
