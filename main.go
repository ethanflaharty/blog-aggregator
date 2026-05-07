package main

import (
	"fmt"

	"github.com/ethanflaharty/blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	err = cfg.SetUser("ethan")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	cfg, err = config.Read()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(cfg)
}
