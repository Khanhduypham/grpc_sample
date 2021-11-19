package main

import (
	"context"
	"example/grpc/calculatorpb"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Println("Sum called...")
	resp := &calculatorpb.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}

	return resp, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PNDRequest,
	stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	log.Println("PND called...")

	k := int32(2)
	n := req.GetNumber()

	for n > 1 {
		if n%k == 0 {
			n /= k
			stream.Send(&calculatorpb.PNDResponse{
				Result: k,
			})
		} else {
			k++
		}
	}
	return nil
}

func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	log.Println("Average called...")
	var total float32
	var count int
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			resp := &calculatorpb.AverageResponse{
				Result: total / float32(count),
			}
			return stream.SendAndClose(resp)
		}
		if err != nil {
			log.Fatalf("Error while Averaga: %v", err)
		}
		total += req.GetNum()
		count++
	}
}

func (*server) FindMax(stream calculatorpb.CalculatorService_FindMaxServer) error {
	log.Println("FindMax called...")
	max := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while findmax: %v", err)
		}
		num := req.GetNum()
		if num > max {
			max = num
		}
		err = stream.Send(&calculatorpb.FindMaxResponse{
			Max: max,
		})
		if err != nil {
			log.Fatalf("Error while sending max %v", err)
			return err
		}
	}
}

func (*server) Sqrt(ctx context.Context, req *calculatorpb.SqrtRequest) (*calculatorpb.SqrtResponse, error) {
	log.Println("Sqrt called...")

	num := req.GetNum()
	if num < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Expect num > 0, num: %v", num)
	}
	resp := &calculatorpb.SqrtResponse{
		Result: math.Sqrt(float64(num)),
	}

	return resp, nil

}

func (*server) SumWithDeadline(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Println("Sum with deadline called...")

	time.Sleep(3 * time.Second)
	if ctx.Err() == context.Canceled {
		log.Println("Client canceled request")
		return nil, status.Errorf(codes.Canceled, "client canceled request")
	}
	resp := &calculatorpb.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}

	return resp, nil
}
func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50069")
	if err != nil {
		log.Fatalf("Error while listening %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	fmt.Println("Server is running")
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Error while serve: %v", err)
	}
}
