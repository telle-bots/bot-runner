package main

import (
	"fmt"
	"time"

	"github.com/telle-bots/bot-runner/pkg/logic"
)

type Button struct {
	Name string `json:"name" name:"Name"`
}

type Keyboard struct {
	ButtonWidth float64   `json:"button_width" name:"Button width"`
	Buttons     []Button  `json:"buttons"      name:"Buttons"`
	Indexes     []int64   `json:"indexes"      name:"Indexes"`
	Points      [][]int64 `json:"points"       name:"Points"`
}

type SendMsg struct {
	ChatID       int64                       `json:"chat_id"      name:"Chat ID"`
	Text         string                      `json:"text"         name:"Text" desc:"Message text to send"`
	Keyboard     *Keyboard                   `json:"keyboard"     name:"Keyboard"`
	Languages    map[string]bool             `json:"languages"    name:"Languages" desc:"Supported languages"`
	Users        map[int64]struct{}          `json:"users"        name:"Users"`
	UserSettings map[int64]map[string]string `json:"userSettings" name:"User settings"`
	SliceOfMaps  []map[bool]string           `json:"sliceOfMaps"  name:"kk"`
}

func main() {
	start := time.Now()
	data := logic.MustStructureOf[SendMsg]()
	fmt.Println(time.Since(start))

	fmt.Println(string(data))
}
