package main

import "blackscout/core"

func main() {
	cfg := core.GetConfig()
	core.Run(cfg)
}
