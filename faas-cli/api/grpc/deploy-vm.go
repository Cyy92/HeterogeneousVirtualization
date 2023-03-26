package grpc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/config"
	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/pb"

	grpcgo "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type DeployVMConfig struct {
	FxGateway string
	FxProcess string

	VMName   string
	Domain   string
	UserData string

	Requests *config.FunctionResources
}

func DeployVM(c DeployVMConfig) error {

	gateway := strings.TrimRight(c.FxGateway, "/")

	conn, err := grpcgo.Dial(gateway, grpcgo.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return errors.New("did not connect: " + err.Error())
	}
	client := pb.NewFxGatewayClient(conn)

	req := pb.CreateVMRequest{
		Instance: c.VMName,
		Domain:   c.Domain,
		UserData: c.UserData,
	}

	hasRequests := false
	req.Requests = &pb.FunctionResources{}
	if c.Requests != nil && len(c.Requests.Memory) > 0 {
		hasRequests = true
		req.Requests.Memory = c.Requests.Memory
	}
	if !hasRequests {
		req.Requests = nil
	}

	_, statusErr := client.DeployVM(context.Background(), &req)
	st, ok := status.FromError(statusErr)
	if !ok {
		return errors.New("Invaild status error.")
	}
	if st.Code() == codes.AlreadyExists {
		fmt.Printf("VM %s already exists. failed deploying.\n", c.VMName)
	}
	if statusErr != nil {
		return statusErr
	}

	return nil
}
