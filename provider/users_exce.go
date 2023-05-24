package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	pb "github.com/travix/gotf-example/pb"
	"github.com/travix/gotf-example/providerpb"
)

var _ providerpb.UsersDataSourceExec = &usersExec{}

type usersExec struct {
	client pb.UserServiceClient
}

func (u *usersExec) Read(ctx context.Context, _ datasource.ReadRequest, _ *datasource.ReadResponse, _ *pb.Users) (*pb.Users, diag.Diagnostics) {
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
