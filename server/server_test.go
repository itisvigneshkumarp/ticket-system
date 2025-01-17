package main

import (
	"context"
	"testing"

	tspb "github.com/itisvigneshkumarp/ticket-system/proto"
	"github.com/stretchr/testify/assert"
)

func TestPurchaseTicket(t *testing.T) {
	server := newServer(50)

	req := &tspb.PurchaseTicketRequest{
		From:      "London",
		To:        "France",
		FirstName: "Alice",
		LastName:  "Smith",
		Email:     "alice.smith@example.com",
		Section:   "A",
	}

	resp, err := server.PurchaseTicket(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "London", resp.From)
	assert.Equal(t, "France", resp.To)
	assert.Equal(t, "Alice", resp.FirstName)
	assert.Equal(t, "Smith", resp.LastName)
	assert.Equal(t, "alice.smith@example.com", resp.Email)
	assert.Equal(t, float32(20.0), resp.PricePaid)
}

func TestGetReceipt(t *testing.T) {
	server := newServer(50)

	req := &tspb.PurchaseTicketRequest{
		From:      "London",
		To:        "France",
		FirstName: "Bob",
		LastName:  "Jones",
		Email:     "bob.jones@example.com",
		Section:   "B",
	}
	resp, err := server.PurchaseTicket(context.Background(), req)
	assert.NoError(t, err)

	receiptReq := &tspb.GetReceiptRequest{
		ReceiptId: resp.ReceiptId,
	}
	receiptResp, err := server.GetReceipt(context.Background(), receiptReq)
	assert.NoError(t, err)
	assert.NotNil(t, receiptResp)
	assert.Equal(t, resp.ReceiptId, receiptResp.ReceiptId)
}

func TestViewUsersBySection(t *testing.T) {
	server := newServer(50)

	req1 := &tspb.PurchaseTicketRequest{
		From:      "London",
		To:        "France",
		FirstName: "Charlie",
		LastName:  "Brown",
		Email:     "charlie.brown@example.com",
		Section:   "A",
	}
	server.PurchaseTicket(context.Background(), req1)

	req2 := &tspb.PurchaseTicketRequest{
		From:      "London",
		To:        "France",
		FirstName: "Dana",
		LastName:  "White",
		Email:     "dana.white@example.com",
		Section:   "B",
	}
	server.PurchaseTicket(context.Background(), req2)

	viewReq := &tspb.ViewUsersBySectionRequest{
		Section: "A",
	}
	viewResp, err := server.ViewUsersBySection(context.Background(), viewReq)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(viewResp.Users))
	assert.Equal(t, "Charlie", viewResp.Users[0].FirstName)
}

func TestRemoveUser(t *testing.T) {
	server := newServer(50)

	req := &tspb.PurchaseTicketRequest{
		From:      "London",
		To:        "France",
		FirstName: "Eve",
		LastName:  "Johnson",
		Email:     "eve.johnson@example.com",
		Section:   "A",
	}
	resp, err := server.PurchaseTicket(context.Background(), req)
	assert.NoError(t, err)

	removeReq := &tspb.RemoveUserRequest{
		ReceiptId: resp.ReceiptId,
	}
	removeResp, err := server.RemoveUser(context.Background(), removeReq)
	assert.NoError(t, err)
	assert.True(t, removeResp.Success)
}

func TestModifySeat(t *testing.T) {
	server := newServer(50)

	req := &tspb.PurchaseTicketRequest{
		From:      "London",
		To:        "France",
		FirstName: "Frank",
		LastName:  "Smith",
		Email:     "frank.smith@example.com",
		Section:   "A",
	}
	resp, err := server.PurchaseTicket(context.Background(), req)
	assert.NoError(t, err)

	modifyReq := &tspb.ModifySeatRequest{
		ReceiptId:  resp.ReceiptId,
		NewSection: "B",
	}
	modifyResp, err := server.ModifySeat(context.Background(), modifyReq)
	assert.NoError(t, err)
	assert.True(t, modifyResp.Success)

	// Verify seat modification
	receiptReq := &tspb.GetReceiptRequest{ReceiptId: resp.ReceiptId}
	receiptResp, err := server.GetReceipt(context.Background(), receiptReq)
	assert.NoError(t, err)
	assert.Equal(t, "B", receiptResp.Section)
}

func TestAnalytics(t *testing.T) {
	server := newServer(50)

	req := &tspb.PurchaseTicketRequest{
		From:      "London",
		To:        "France",
		FirstName: "Frank",
		LastName:  "Smith",
		Email:     "frank.smith@example.com",
		Section:   "A",
	}
	_, err := server.PurchaseTicket(context.Background(), req)
	assert.NoError(t, err)

	analyticsReq := &tspb.GetAnalyticsRequest{}
	analyticsResp, err := server.GetAnalytics(context.Background(), analyticsReq)
	assert.NoError(t, err)
	assert.Equal(t, int32(100), analyticsResp.TotalTickets)
	assert.Equal(t, int32(1), analyticsResp.TotalTicketsSold)
	assert.Equal(t, int32(99), analyticsResp.TotalTicketsAvailable)
	assert.Equal(t, float64(20.0), analyticsResp.TotalRevenueInDollars)
	assert.Equal(t, int32(1), analyticsResp.SectionAOccupancy)
}
