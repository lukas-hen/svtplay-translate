package main

func GetKeysFromMap[A comparable, B any](some_map map[A]B) []A {

	out := make([]A, len(some_map))

	i := 0
	for k := range some_map {
		out[i] = k
		i++
	}

	return out
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
