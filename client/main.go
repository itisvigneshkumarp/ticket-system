package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/itisvigneshkumarp/ticket-system/pkg/util"
	tspb "github.com/itisvigneshkumarp/ticket-system/proto"
)

func main() {
	serverURL := flag.String("server", "localhost:50051", "gRPC Server URL")
	flag.Parse()

	// Connect to the server
	conn, err := grpc.NewClient(*serverURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := tspb.NewTrainTicketServiceClient(conn)
	ctx := context.Background()

	// Purchase a ticket
	purchaseRequest := &tspb.PurchaseTicketRequest{
		From:      "London",
		To:        "France",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Section:   "A",
	}
	purchaseResponse, err := client.PurchaseTicket(ctx, purchaseRequest)
	if err != nil {
		log.Fatalf("Error purchasing ticket: %v", err)
	}
	fmt.Printf("Ticket Purchased: %+v\n", util.ToJSON(purchaseResponse))

	// Get Receipt
	receiptRequest := &tspb.GetReceiptRequest{
		ReceiptId: purchaseResponse.ReceiptId,
	}
	receiptResponse, err := client.GetReceipt(ctx, receiptRequest)
	if err != nil {
		log.Fatalf("Error getting receipt: %v", err)
	}
	fmt.Printf("Receipt Details: %+v\n", util.ToJSON(receiptResponse))

	// Analytics
	analyticsRequest := &tspb.GetAnalyticsRequest{}
	analyticsResponse, err := client.GetAnalytics(ctx, analyticsRequest)
	if err != nil {
		log.Fatalf("Error fetching analytics: %v", err)
	}
	fmt.Printf("Analytics: %+v\n", util.ToJSON(analyticsResponse))

	// View Users by Section
	viewRequest := &tspb.ViewUsersBySectionRequest{
		Section: "A",
	}
	viewResponse, err := client.ViewUsersBySection(ctx, viewRequest)
	if err != nil {
		log.Fatalf("Error viewing users by section: %v", err)
	}
	fmt.Printf("Users in Section A: %+v\n", util.ToJSON(viewResponse.Users))

	// Modify Seat
	modifyRequest := &tspb.ModifySeatRequest{
		ReceiptId:  purchaseResponse.ReceiptId,
		NewSection: "B",
	}
	modifyResponse, err := client.ModifySeat(ctx, modifyRequest)
	if err != nil {
		log.Fatalf("Error modifying seat: %v", err)
	}
	fmt.Printf("Seat Modification: %+v\n", util.ToJSON(modifyResponse))

	// Remove User
	removeRequest := &tspb.RemoveUserRequest{
		ReceiptId: purchaseResponse.ReceiptId,
	}
	removeResponse, err := client.RemoveUser(ctx, removeRequest)
	if err != nil {
		log.Fatalf("Error removing user: %v", err)
	}
	fmt.Printf("User Removal: %+v\n", util.ToJSON(removeResponse))
}
