package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type EnvPair struct {
	key   string
	value interface{}
}

func main() {
	pairs, err := getPairs()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	res, err := getResult(pairs)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	js, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(string(js))
}

func getResult(pairs []EnvPair) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	for i := range pairs {
		// For every _ in the name we assume an additional level of nesting in the resulting json
		keys := strings.Split(pairs[i].key, "_")
		currentMap := res
		for j := 0; j < len(keys); j++ {
			// If this is the last key after the split,
			// then set the value instead of creating a new object
			if len(keys)-1 == j {
				currentMap[keys[j]] = pairs[i].value
				// fmt.Printf("Done for %+v\n", currentMap)
				continue
			}

			// New JSON object
			var newMap map[string]interface{}
			if _, exists := currentMap[keys[j]]; exists {
				// If there already exists an object for this key, then add to it
				tmpMap, ok := currentMap[keys[j]].(map[string]interface{})
				if !ok {
					// return nil, fmt.Errorf("Unexpected map structure for %+v", currentMap[keys[j]])
					continue
				}
				newMap = tmpMap
			} else {
				// or else create a new map/object
				newMap = make(map[string]interface{})
			}
			// For our current key list after the split we can work on this new object
			currentMap[keys[j]] = newMap
			currentMap = newMap
		}
	}
	return res, nil
}

func getPairs() ([]EnvPair, error) {
	pairs := make([]EnvPair, len(os.Environ()))
	for i := range os.Environ() {
		pair := strings.SplitN(os.Environ()[i], "=", 2)
		if len(pair) != 2 {
			// return nil, fmt.Errorf("Env in wrong format! Got %v", pair)
			// fmt.Printf("Env in wrong format! Got %v", pair)
			continue
		}
		vals := strings.Split(pair[1], ":")
		if len(vals) == 1 {
			pairs[i] = EnvPair{pair[0], vals[0]}
			continue
		}
		pairs[i] = EnvPair{pair[0], vals}
	}
	return pairs, nil
}
