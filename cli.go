package main

import (
	"fmt"
	"strconv"
)

func PromptUserDecision[T comparable](opts []T, prompt string) T {

	fmt.Println(prompt)

	for idx, opt := range opts {
		// Unsafe cast of dynamic type to string here.
		fmt.Printf("\t%d. %s\n", idx, opt)
	}

	fmt.Printf("> ")
	var choice string
	_, err := fmt.Scanf("%s", &choice)
	check(err)

	intChoice, err := strconv.Atoi(choice)
	check(err)
	if intChoice < 0 || intChoice > len(opts) {
		panic("Selected choice was not in range of the given options.")
	}

	return opts[intChoice]
}
