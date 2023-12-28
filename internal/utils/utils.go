package utils

import (
	"fmt"
	"strconv"
)

func GetKeysFromMap[A comparable, B any](some_map map[A]B) []A {

	out := make([]A, len(some_map))

	i := 0
	for k := range some_map {
		out[i] = k
		i++
	}

	return out
}

func PromptUserDecision[T string](opts []T, prompt string) T {

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

func check(err error) {
	if err != nil {
		panic(err)
	}
}
