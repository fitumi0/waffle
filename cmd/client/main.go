package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/fitumi0/waffle/core/client"
	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()
	client, err := client.NewClient(ctx, "localhost:50051", uuid.New().String())
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		if err := client.SendMessage("test", line); err != nil {
			fmt.Fprintln(os.Stderr, "Error sending message:", err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading stdin:", err)
	}
}
