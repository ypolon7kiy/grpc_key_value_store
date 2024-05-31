package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	pb "client/kvstore"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = "localhost:50051"
	}

	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewKeyValueStoreClient(conn)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter command (set|get|delete|exit) k [v]: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("failed to read input: %v", err)
		}

		input = strings.TrimSpace(input)
		if input == "exit" {
			fmt.Println("Exiting...")
			break
		}

		handleCommand(client, input)
	}
}

func handleCommand(client pb.KeyValueStoreClient, input string) {
	parts := strings.Split(input, " ")
	if len(parts) < 2 {
		fmt.Println("Invalid command. Usage: set <key> <value>, get <key>, delete <key>")
		return
	}

	cmd := parts[0]
	key := parts[1]
	var value string
	if len(parts) > 2 {
		value = strings.Join(parts[2:], " ")
	}

	switch cmd {
	case "set":
		if value == "" {
			fmt.Println("Usage for set: set <key> <value>")
			return
		}
		setResponse, err := client.Set(context.Background(), &pb.SetRequest{Key: key, Value: value})
		if err != nil {
			log.Printf("could not set value: %v", err)
			return
		}
		log.Printf("Set Response: %s", setResponse.Message)

	case "get":
		getResponse, err := client.Get(context.Background(), &pb.GetRequest{Key: key})
		if err != nil {
			log.Printf("could not get value: %v", err)
			return
		}
		log.Printf("Get Response: Value: %s, Message: %s", getResponse.Value, getResponse.Message)

	case "delete":
		deleteResponse, err := client.Delete(context.Background(), &pb.DeleteRequest{Key: key})
		if err != nil {
			log.Printf("could not delete key: %v", err)
			return
		}
		log.Printf("Delete Response: %s", deleteResponse.Message)

	default:
		fmt.Println("Invalid command. Usage: set <key> <value>, get <key>, delete <key>")
	}
}
