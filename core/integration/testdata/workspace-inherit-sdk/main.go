package main

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()
	client, err := dagger.Connect(ctx)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	contents, err := client.CurrentWorkspace().
		File("/core/integration/testdata/workspace-inherit-sdk/outer-marker.txt").
		Contents(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Print(contents)
}
