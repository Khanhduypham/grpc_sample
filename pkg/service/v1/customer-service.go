package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"

	pb "example/grpc/pkg/api/v1"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

const (
	port        = ":50051"
	DB_USER     = "postgres"
	DB_PASSWORD = "nelson"
	DB_NAME     = "grpc"
)

// server is used to implement customer.CustomerServer.
type server struct {
	savedCustomers []*pb.CustomerRequest
	pb.UnimplementedCustomerServer
	db *sql.DB
}

func (s *server) CreateCustomer(ctx context.Context, in *pb.CustomerRequest) (*pb.CustomerResponse, error) {
	s.savedCustomers = append(s.savedCustomers, in)
	return &pb.CustomerResponse{Id: in.Id, Success: true}, nil

}

func (s *server) GetCustomers(filter *pb.CustomerFilter, stream pb.Customer_GetCustomersServer) error {
	for _, customer := range s.savedCustomers {
		if filter.Keyword != "" {
			if !strings.Contains(customer.Name, filter.Keyword) {
				continue
			}
		}
		if err := stream.Send(customer); err != nil {
			return err
		}
	}
	return nil

	/* //get from db
	rows, _ := s.db.Query(`SELECT * FROM public."todo"`)
	for rows.Next() {
		var movieID int32
		var movieName string
		var b string

		rows.Scan(&movieID, &movieName, &b)

		fmt.Printf("%v", movieName)
	}
	*/
}

func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		panic(err)
	}

	return db
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCustomerServer(s, &server{db: setupDB()})
	s.Serve(lis)
}
