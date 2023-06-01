package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	pb "github.com/travix/gotf-example/pb"
	"github.com/travix/gotf-example/provider/providerpb"
)

var _ providerpb.UsersDataSourceExec = &usersExec{}

type usersExec struct {
	client pb.UserServiceClient
}

func (u *usersExec) Read(ctx context.Context, _ datasource.ReadRequest, _ *datasource.ReadResponse, _ *pb.Users) (*pb.Users, diag.Diagnostics) {
	tflog.Info(ctx, "reading users from grpc")
	users, err := u.client.ListUsers(ctx, &pb.Empty{})
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("failed to list users", err.Error())
		return nil, diags
	}
	return users, nil
}

func (u *usersExec) SetUserServiceClient(client pb.UserServiceClient) {
	u.client = client
}
