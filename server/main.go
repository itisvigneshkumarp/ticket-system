package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"strconv"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	tspb "github.com/itisvigneshkumarp/ticket-system/proto"
)

type server struct {
	tspb.UnimplementedTrainTicketServiceServer
	mu          sync.Mutex
	users       map[string]*tspb.PurchaseTicketResponse
	seatsA      map[string]bool
	seatsB      map[string]bool
	counter     int
	sectionSize int
}

func newServer(sectionSize int) *server {
	return &server{
		users:       make(map[string]*tspb.PurchaseTicketResponse),
		seatsA:      make(map[string]bool),
		seatsB:      make(map[string]bool),
		counter:     0,
		sectionSize: sectionSize,
	}
}

func (s *server) allocateSeat(section string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var seats map[string]bool
	if section == "A" {
		seats = s.seatsA
	} else if section == "B" {
		seats = s.seatsB
	} else {
		return "", fmt.Errorf("invalid section")
	}

	for i := 1; i <= s.sectionSize; i++ {
		seat := fmt.Sprintf("%s-%d", section, i)
		if !seats[seat] {
			seats[seat] = true
			return seat, nil
		}
	}

	return "", fmt.Errorf("no seats available in section %s", section)
}

func (s *server) PurchaseTicket(ctx context.Context, req *tspb.PurchaseTicketRequest) (*tspb.PurchaseTicketResponse, error) {
	seat, err := s.allocateSeat(req.Section)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	receiptID := fmt.Sprintf("R-%d", s.counter)
	response := &tspb.PurchaseTicketResponse{
		ReceiptId: receiptID,
		From:      req.From,
		To:        req.To,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		PricePaid: 20.0,
		Seat:      seat,
		Section:   req.Section,
	}

	s.users[receiptID] = response
	return response, nil
}

func (s *server) GetReceipt(ctx context.Context, req *tspb.GetReceiptRequest) (*tspb.GetReceiptResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	response, ok := s.users[req.ReceiptId]
	if !ok {
		return nil, fmt.Errorf("receipt not found")
	}

	return &tspb.GetReceiptResponse{
		ReceiptId: response.ReceiptId,
		From:      response.From,
		To:        response.To,
		FirstName: response.FirstName,
		LastName:  response.LastName,
		Email:     response.Email,
		PricePaid: response.PricePaid,
		Seat:      response.Seat,
		Section:   response.Section,
	}, nil
}

func (s *server) ViewUsersBySection(ctx context.Context, req *tspb.ViewUsersBySectionRequest) (*tspb.ViewUsersBySectionResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var users []*tspb.UserSeatAllocation
	for _, user := range s.users {
		if user.Section == req.Section {
			users = append(users, &tspb.UserSeatAllocation{
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
				Seat:      user.Seat,
			})
		}
	}

	return &tspb.ViewUsersBySectionResponse{Users: users}, nil
}

func (s *server) RemoveUser(ctx context.Context, req *tspb.RemoveUserRequest) (*tspb.RemoveUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[req.ReceiptId]
	if !ok {
		return &tspb.RemoveUserResponse{Success: false}, fmt.Errorf("user not found")
	}

	delete(s.users, req.ReceiptId)
	if user.Section == "A" {
		delete(s.seatsA, user.Seat)
	} else {
		delete(s.seatsB, user.Seat)
	}

	return &tspb.RemoveUserResponse{Success: true}, nil
}

func (s *server) ModifySeat(ctx context.Context, req *tspb.ModifySeatRequest) (*tspb.ModifySeatResponse, error) {
	// Check if the user exists
	s.mu.Lock()
	user, ok := s.users[req.ReceiptId]
	if !ok {
		s.mu.Unlock()
		return &tspb.ModifySeatResponse{Success: false}, fmt.Errorf("user not found")
	}
	oldSeat := user.Seat
	oldSection := user.Section
	s.mu.Unlock()

	// Allocate a new seat outside the mutex to avoid deadlocks
	newSeat, err := s.allocateSeat(req.NewSection)
	if err != nil {
		return &tspb.ModifySeatResponse{Success: false}, err
	}

	// Update user details with a lock
	s.mu.Lock()
	defer s.mu.Unlock()

	// Release the old seat
	if oldSection == "A" {
		delete(s.seatsA, oldSeat)
	} else {
		delete(s.seatsB, oldSeat)
	}

	// Update user information
	user.Seat = newSeat
	user.Section = req.NewSection
	return &tspb.ModifySeatResponse{Success: true}, nil
}

func (s *server) GetAnalytics(ctx context.Context, req *tspb.GetAnalyticsRequest) (*tspb.GetAnalyticsResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	totalTicketsSold := len(s.users)
	totalRevenue := float64(totalTicketsSold) * 20.0
	sectionAOccupancy := len(s.seatsA)
	sectionBOccupancy := len(s.seatsB)

	return &tspb.GetAnalyticsResponse{
		TotalTickets:          int32(2 * s.sectionSize),
		TotalTicketsSold:      int32(totalTicketsSold),
		TotalTicketsAvailable: int32((2 * s.sectionSize) - (sectionAOccupancy + sectionBOccupancy)),
		TotalRevenue:          totalRevenue,
		SectionAOccupancy:     int32(sectionAOccupancy),
		SectionBOccupancy:     int32(sectionBOccupancy),
	}, nil
}

func main() {
	sectionSizeArg := flag.String("seats", "50", "Number of seats per section")
	flag.Parse()

	sectionSize, err := strconv.Atoi(*sectionSizeArg)
	if err != nil {
		panic(fmt.Errorf("invalid section size: %v", err))
	}

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	tspb.RegisterTrainTicketServiceServer(grpcServer, newServer(sectionSize))
	reflection.Register(grpcServer)

	fmt.Println("Server is running on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
