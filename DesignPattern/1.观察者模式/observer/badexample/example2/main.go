package main

import (
	"design-pattern/observer/badexample/example2/player"
)

func main() {
	p := player.NewPlayer("xie4ever")
	_ = p.Attack()
}
