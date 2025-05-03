package api

import (
	"regexp"
	"fmt"
	"strconv"
)

func getToken(message string) uint64 {
	fmt.Println("DEBUG - Authorization header:", message) // DEBUG
	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
	stringToken := re.FindAllString(message, -1)

	// DEBUG
	if len(stringToken) == 0 {
        fmt.Println("DEBUG - no token found!") 
        return 0
    }

	// DEBUG
	token, err := strconv.ParseUint(stringToken[0], 10, 64)
    if err != nil {
        fmt.Printf("DEBUG - invalid token %q: %v\n", stringToken[0], err)
        return 0
    }

	// token, _ := strconv.Atoi(stringToken[0])
	return uint64(token)
}