package grpc

import (
	"context"
	"errors"
	"strings"

	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/pb"
	grpcgo "google.golang.org/grpc"
)

func DeleteVM(fxGateway, vmName string) error {

	gateway := strings.TrimRight(fxGateway, "/")

	conn, err := grpcgo.Dial(gateway, grpcgo.WithInsecure())
	if err != nil {
		return errors.New("did not connect: " + err.Error())
	}
	client := pb.NewFxGatewayClient(conn)

	_, statusErr := client.DeleteVM(context.Background(), &pb.DeleteVMRequest{VMName: vmName})
	if statusErr != nil {
		return errors.New("did not delete: " + statusErr.Error())
	}

	return nil
}
