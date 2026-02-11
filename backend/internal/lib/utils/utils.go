package utils

import (
	"encoding/json"
	"fmt"
)

func PrintJSON(v interface{}) {
	json, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	fmt.Println("JSON:", string(json))
}

// we will define utils for background tasks , in this lib directory , these tasks is used for delegating tasks asynchronously , we will be using asynq library for
// this purpose e.g. sending notification , sending email.
// One of the solution is goroutine. We can delegate the task on one goroutine let it run concurrently, and when finished just notify the things .

// -- External library used in this folder - encoding/json, fmt
