package commands

import (
	"fmt"
	"strconv"
)

func extractServerIDs(s []string) ([]int, error) {
	serverIDs := []int{}

	for _, e := range s {
		i, err := strconv.Atoi(e)
		if err != nil {
			return nil, fmt.Errorf("Provided value [%v] for server id is not of type int", e)
		}
		serverIDs = append(serverIDs, i)
	}

	return serverIDs, nil
}
