package main

import (
	"fmt"
	"lgo"
	"lgo/log"
)

type Person struct {
	Name string
	Age  int
}

func init() {
	lgo.HandleFunc("/getAdLineDataAPI", AdLineInfo)
}

func AdLineInfo(ctx *lgo.Context) {
	fmt.Println("adLineInfo")
	person := &Person{"leroy", 20}

	log.Debug("select * from ")
	defer ctx.SetResultBody(person)
}
