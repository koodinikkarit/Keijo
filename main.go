package main

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"

	"github.com/chbmuc/cec"

	KeijoService "github.com/koodinikkarit/keijo/service"
)

type server struct {
	c *cec.Connection
}

func (s *server) TurnOn(ctx context.Context, in *KeijoService.TurnOnRequest) (*KeijoService.TurnOnResponse, error) {
	fmt.Println("Turn on ", int(in.Address))
	s.c.PowerOn(int(in.Address))
	return &KeijoService.TurnOnResponse{}, nil
}

func (s *server) TurnOff(ctx context.Context, in *KeijoService.TurnOffRequest) (*KeijoService.TurnOffResponse, error) {
	s.c.Standby(int(in.Address))
	return &KeijoService.TurnOffResponse{}, nil
}

func (s *server) ChangeSource(ctx context.Context, in *KeijoService.ChangeSourceRequest) (*KeijoService.ChangeSourceResponse, error) {
	var sourceChar string
	var destinationChar string
	source := in.Source
	if source > 9 {
		switch source {
		case 10:
			sourceChar = "A"
		case 11:
			sourceChar = "B"
		case 12:
			sourceChar = "C"
		case 13:
			sourceChar = "D"
		case 14:
			sourceChar = "E"
		case 15:
			sourceChar = "F"
		}
	} else {
		sourceChar = strconv.Itoa(int(source))
	}

	if in.Destination > 9 {
		switch in.Destination {
		case 10:
			destinationChar = "A"
		case 11:
			destinationChar = "B"
		case 12:
			destinationChar = "C"
		case 13:
			destinationChar = "D"
		case 14:
			destinationChar = "E"
		case 15:
			destinationChar = "F"
		}
	} else {
		destinationChar = strconv.Itoa(int(in.Destination))
	}
	s.c.Transmit(sourceChar + destinationChar + ":82:" + strconv.Itoa(int(in.SourceNumber)) + "0:00")
	return &KeijoService.ChangeSourceResponse{}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":12345")
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	c, err := cec.Open("", "cec.go")
	if err != nil {
		fmt.Println(err)
	}

	KeijoService.RegisterKeijoServer(s, &server{c: c})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
