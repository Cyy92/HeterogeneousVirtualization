package grpc

import (
	"context"
	"errors"
	"strings"

	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/pb"
	grpcgo "google.golang.org/grpc"
)

func InvokeVM(VMName, fxGateway string, input []byte) (string, error) {

	gateway := strings.TrimRight(fxGateway, "/")

	conn, err := grpcgo.Dial(gateway, grpcgo.WithInsecure())
	if err != nil {
		return "", errors.New("did not connect: " + err.Error())
	}
	client := pb.NewFxGatewayClient(conn)

	message, statusErr := client.InvokeVM(context.Background(), &pb.InvokeServiceRequest{Service: VMName, Input: input})
	if statusErr != nil {
		return "", errors.New("did not invoke: " + statusErr.Error())
	}

	return message.Msg, nil
}
