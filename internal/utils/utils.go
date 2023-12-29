package utils

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

// Returns all the keys from a map in a list.
func GetKeysFromMap[A comparable, B any](some_map map[A]B) []A {

	out := make([]A, len(some_map))

	i := 0
	for k := range some_map {
		out[i] = k
		i++
	}

	return out
}

// Takes a slice containing type T and prompts the user to choose one element with the prompt string.
// Returns the option of type T that the user selects.
func PromptUserDecision[T string](opts []T, prompt string) T {

	fmt.Println(prompt)

	for idx, opt := range opts {
		fmt.Printf("\t%d. %s\n", idx, opt)
	}

	fmt.Printf("> ")
	var choice string
	_, err := fmt.Scanf("%s", &choice)
	if err != nil {
		log.Fatalln("Failed to scan user input.")
	}

	intChoice, err := strconv.Atoi(choice)
	if err != nil {
		log.Fatalf("Could not convert choice into an integer.\n")
	}

	if intChoice < 0 || intChoice > len(opts) {
		log.Fatalln("Selected choice was not in range of the given options.")
	}

	return opts[intChoice]
}

// Taken from: https://gist.github.com/schwarzeni/f25031a3123f895ff3785970921e962c
func GetInterfaceIpv4Addr(interfaceName string) (addr string, err error) {
	var (
		ief      *net.Interface
		addrs    []net.Addr
		ipv4Addr net.IP
	)
	if ief, err = net.InterfaceByName(interfaceName); err != nil { // get interface
		return
	}
	if addrs, err = ief.Addrs(); err != nil { // get addresses
		return
	}
	for _, addr := range addrs { // get ipv4 address
		if ipv4Addr = addr.(*net.IPNet).IP.To4(); ipv4Addr != nil {
			break
		}
	}
	if ipv4Addr == nil {
		return "", fmt.Errorf(fmt.Sprintf("interface %s don't have an ipv4 address\n", interfaceName))
	}
	return ipv4Addr.String(), nil
}
