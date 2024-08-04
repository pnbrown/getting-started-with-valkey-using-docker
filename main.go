package main

import (
	"context"
	"fmt"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/valkey-io/valkey-go"
)

const (
	keyData   = "new-redis"
	valueData = "Valkey"
)

func main() {
	ctx := context.Background()

	containerRequest := testcontainers.ContainerRequest{
        Name:         "valkey",
		Image:        "valkey/valkey:7.2.5",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerRequest,
		Started:          true,
		Reuse:            true,
	})
	if err != nil {
		log.Fatalf("Could not start Valkey: %s", err)
	}
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			log.Fatalf("Unable to stop Valkey: %s", err)
		}
	}()

	endpoint, err := container.Endpoint(ctx, "")
	if err != nil {
		log.Fatalf("Unable to retrieve the endpoint: %s", err)
	}

	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{endpoint},
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	err = client.Do(ctx, client.B().Set().Key(keyData).Value(valueData).Build()).Error()
	if err != nil {
		panic(err)
	}

	value, err := client.Do(ctx, client.B().Get().Key(keyData).Build()).ToString()
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
}

