package main

import "fmt"

func main() {
	input := make(chan string)

	go func(i chan string) {
		for {
			fmt.Println(<-i)
		}

	}(input)

	w := newUi(input)
	(*w).ShowAndRun()
}
