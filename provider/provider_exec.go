package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_grpc "github.com/travix/gotf-example/example-server/grpc"
	"github.com/travix/gotf-example/pb"
	"github.com/travix/gotf-example/provider/providerpb"
)

var _ providerpb.ExampleExec = &ProviderExec{}

type ProviderExec struct {
}

func (p *ProviderExec) ConfigureGrpc(ctx context.Context, model *pb.ProviderModel) (conn grpc.ClientConnInterface, diagnostics diag.Diagnostics) {
	// credentials and serverAddr can be fetched from req.Config by setting
	opts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(_grpc.NewClientAuthInterceptor(model.KeyId, model.SecretKey)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	var err error
	if model.Endpoint == "" {
		model.Endpoint = "127.0.0.1:50051"
	}
	tflog.Info(ctx, fmt.Sprintf("dialing grpc connection with example grcp '%s'", model.Endpoint))
	conn, err = grpc.Dial(model.Endpoint, opts...)
	if err != nil {
		diagnostics.AddError("Failed connecting to example grcp", fmt.Sprintf("eror in grpc connection with %s: %v", model.Endpoint, err))
		return nil, diagnostics
	}
	return
}

func (p *ProviderExec) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		providerpb.NewUsersDataSource(&usersExec{}),
		providerpb.NewGroupsDataSource(&groupsExec{}),
	}
}

func (p *ProviderExec) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		providerpb.NewUserResource(&userExec{}),
		providerpb.NewGroupResource(&groupExec{}),
	}
}
