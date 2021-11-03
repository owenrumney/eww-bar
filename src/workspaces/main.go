package main

// NICKED FROM https://github.com/liamg/dotfiles/blob/master/eww/src/workspaces/main.go

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"go.i3wm.org/i3/v4"
)

func main() {

	if len(os.Args) <= 1 {
		panic("please supply output number")
	}

	index, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	updateWorkspaces(index)

	subscription := i3.Subscribe(i3.WorkspaceEventType)
	for subscription.Next() {
		event := subscription.Event()
		switch event.(type) {
		case *i3.WorkspaceEvent:
			updateWorkspaces(index)
		}
	}
}

func updateWorkspaces(outputIndex int) {
	workspaces, err := i3.GetWorkspaces()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	var outputName string
	outputs, err := i3.GetOutputs()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	if outputIndex < len(outputs) {
		outputName = outputs[outputIndex].Name
	}

	// open box
	fmt.Printf(`(box 	:orientation "h"	:space-evenly false  :spacing 10`)
	for _, workspace := range workspaces {

		if workspace.Output != outputName {
			continue
		}

		var class string
		if workspace.Urgent {
			class = "urgent"
		} else if workspace.Focused {
			class = "focused"
		}

		fmt.Printf(
			`(button `+
				`:onclick "i3-msg workspace '%s'"`+
				`:class '%s'`+
				`(label :text '%s'))`,
			workspace.Name,
			class,
			strings.TrimPrefix(workspace.Name, strconv.Itoa(int(workspace.Num))),
		)
	}

	// close box + newline for send
	fmt.Println(")")
}
