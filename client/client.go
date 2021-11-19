package main

import (
	"context"
	"example/grpc/calculatorpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	cc, err := grpc.Dial("localhost:50069", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Error while dial: %v", err)
	}

	defer cc.Close()

	client := calculatorpb.NewCalculatorServiceClient(cc)

	//callSum(client)
	//callPND(client)
	//callAverage(client)
	//callFindMax(client)
	//callSqrt(client)
	callSumWithDeadline(client, 2*time.Second)
}

func callSum(c calculatorpb.CalculatorServiceClient) {
	res, err := c.Sum(context.Background(), &calculatorpb.SumRequest{
		Num1: 5,
		Num2: 6,
	})

	if err != nil {
		log.Fatalf("Error while call sum at client: %v", err)
	}

	log.Printf("Sum response: %v", res.GetResult())
}

func callPND(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.PrimeNumberDecomposition(context.Background(), &calculatorpb.PNDRequest{
		Number: 120,
	})

	if err != nil {
		log.Fatalf("Error while call PND at client %v", err)
	}

	for {
		res, receiveErr := stream.Recv()
		if receiveErr == io.EOF {
			log.Fatalln("Server finished streaming")
			return
		}

		log.Printf("Prime number response: %v", res.GetResult())
	}

}

func callAverage(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("Error when getting stream: %v", err)
	}

	list := []calculatorpb.AverageRequest{
		calculatorpb.AverageRequest{
			Num: 1.9,
		},
		calculatorpb.AverageRequest{
			Num: 1.9,
		},
		calculatorpb.AverageRequest{
			Num: 1.9,
		},
		calculatorpb.AverageRequest{
			Num: 1.9,
		},
		calculatorpb.AverageRequest{
			Num: 5,
		},
	}

	for _, req := range list {
		err := stream.Send(&req)
		if err != nil {
			log.Fatalf("Error while send request %v", err)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error when getting response %v", err)
	}
	log.Printf("Average: %v", resp)
}

func callFindMax(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.FindMax(context.Background())
	if err != nil {
		log.Fatalf("Error when getting stream: %v", err)
	}

	list := []calculatorpb.FindMaxRequest{
		calculatorpb.FindMaxRequest{
			Num: 12,
		},
		calculatorpb.FindMaxRequest{
			Num: 3,
		},
		calculatorpb.FindMaxRequest{
			Num: 9,
		},
		calculatorpb.FindMaxRequest{
			Num: 20,
		},
	}

	waitc := make(chan struct{})
	go func() {
		for _, req := range list {
			err := stream.Send(&req)
			if err != nil {
				log.Fatalf("Error while send request %v", err)
			}
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				log.Fatalln("Server finished streaming")
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving result from server: %v", err)
				break
			}
			log.Printf("Current Max: %v", resp.GetMax())
		}
		close(waitc)
	}()

	<-waitc
}

func callSqrt(c calculatorpb.CalculatorServiceClient) {
	res, err := c.Sqrt(context.Background(), &calculatorpb.SqrtRequest{
		Num: -16,
	})

	if err != nil {
		log.Printf("Error while call sqrt: %v", err)
		if errStatus, ok := status.FromError(err); ok {
			log.Printf("err msg: %v\n", errStatus.Message())
			log.Printf("err code: %v\n", errStatus.Code())
		}

	} else {
		log.Printf("Sqrt response: %v", res.GetResult())
	}
}

func callSumWithDeadline(c calculatorpb.CalculatorServiceClient, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := c.SumWithDeadline(ctx, &calculatorpb.SumRequest{
		Num1: 5,
		Num2: 6,
	})

	if err != nil {
		if statusErr, ok := status.FromError(err); ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Println("Deadline exceeded")
			} else {
				log.Printf("Error while calling sum dl: %v", err)
			}
		} else {
			log.Fatalf("Error while call sum with deadline (unknown): %v", err)
		}
	}

	log.Printf("Sum response: %v", res.GetResult())
}
