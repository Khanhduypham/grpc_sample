package main

import (
	"context"
	"encoding/json"
	pb "example/grpc/pkg/api/v1"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

// createCustomer calls the RPC method CreateCustomer of CustomerServer
func createCustomer(client pb.CustomerClient, customer *pb.CustomerRequest) {

	resp, err := client.CreateCustomer(context.Background(), customer)
	if err != nil {
		log.Fatalf("Could not create Customer: %v", err)
	}
	if resp.Success {
		log.Printf("A new Customer has been added with id: %d", resp.Id)
	}
}

// getCustomers calls the RPC method GetCustomers of CustomerServer
func getCustomers(client pb.CustomerClient, filter *pb.CustomerFilter) {

	// calling the streaming API
	stream, err := client.GetCustomers(context.Background(), filter)
	if err != nil {
		log.Fatalf("Error on get customers: %v", err)
	}
	for {
		customer, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetCustomers(_) = _, %v", client, err)
		}
		log.Printf("Customer: %v", customer)
	}
}

func GetCustomerHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Creates a new CustomerClient
	client := pb.NewCustomerClient(conn)

	customer := &pb.CustomerRequest{
		Id:    1,
		Name:  "CongPV 1",
		Email: "vancogn1@gmail.com",
		Phone: "123456789",
		Addresses: []*pb.CustomerRequest_Address{
			&pb.CustomerRequest_Address{
				Street:            "111C Nguyen Lam",
				City:              "TPHCM",
				State:             "TP",
				Zip:               "124",
				IsShippingAddress: false,
			},
			&pb.CustomerRequest_Address{
				Street:            "111B Nguyen Lam",
				City:              "TPHCM",
				State:             "TP",
				Zip:               "124",
				IsShippingAddress: true,
			},
		},
	}

	// Create a new customer
	createCustomer(client, customer)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("")

	//fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/customers", GetCustomerHandler).Methods("GET")
	//r.HandleFunc("/customers", CreateBook).Methods("POST")
	http.ListenAndServe(":8080", r)
	// Set up a connection to the gRPC server.

}
