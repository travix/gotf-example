package main

import (
	"context"
	"github.com/travix/gotf-example/provider/providerpb"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/travix/gotf-example/pb"
)

var _ providerpb.UserResourceExec = &userExec{}

type userExec struct {
	client pb.UserServiceClient
}

func (u *userExec) SetUserServiceClient(client pb.UserServiceClient) {
	u.client = client
}

func (u *userExec) Create(ctx context.Context, _ resource.CreateRequest, _ *resource.CreateResponse, data *pb.User) (*pb.User, diag.Diagnostics) {
	user, err := u.client.CreateUser(ctx, data)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("failed to create user", err.Error())
		return nil, diags
	}
	return user, nil
}

func (u *userExec) Read(ctx context.Context, _ resource.ReadRequest, _ *resource.ReadResponse, data *pb.User) (*pb.User, diag.Diagnostics) {
	user, err := u.client.GetUser(ctx, &pb.GetUserRequest{Username: data.Username})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return data, nil
		}
		var diags diag.Diagnostics
		diags.AddError("failed to get user", err.Error())
		return nil, diags
	}
	return user, nil
}

func (u *userExec) Update(ctx context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse, data *pb.User) (*pb.User, diag.Diagnostics) {
	user, err := u.client.UpdateUser(ctx, data)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("failed to update user", err.Error())
		return nil, diags
	}
	return user, nil
}

func (u *userExec) Delete(ctx context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse, data *pb.User) diag.Diagnostics {
	_, err := u.client.DeleteUser(ctx, data)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("failed to delete user", err.Error())
		return diags
	}
	return nil
}
