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

var _ providerpb.GroupResourceExec = &groupExec{}

type groupExec struct {
	groupClient pb.GroupServiceClient
	userClient  pb.UserServiceClient
}

func (g *groupExec) Create(ctx context.Context, _ resource.CreateRequest, _ *resource.CreateResponse, data *pb.Group) (*pb.Group, diag.Diagnostics) {
	user, err := g.groupClient.CreateGroup(ctx, data)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("failed to create group", err.Error())
		return nil, diags
	}
	if len(data.Users) > 0 {
		for _, user := range data.Users {
			_, err := g.userClient.CreateUser(ctx, user)
			if err != nil {
				var diags diag.Diagnostics
				diags.AddError("failed to add users from group", err.Error())
				return nil, diags
			}
		}
	}
	return user, nil
}

func (g *groupExec) Read(ctx context.Context, _ resource.ReadRequest, _ *resource.ReadResponse, data *pb.Group) (*pb.Group, diag.Diagnostics) {
	group, err := g.groupClient.GetGroup(ctx, &pb.GetGroupRequest{Name: data.Name})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return data, nil
		}
		var diags diag.Diagnostics
		diags.AddError("failed to get group", err.Error())
		return nil, diags
	}
	return group, nil
}

func (g *groupExec) Update(ctx context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse, data *pb.Group) (*pb.Group, diag.Diagnostics) {
	group, err := g.groupClient.UpdateGroup(ctx, data)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("failed to update group", err.Error())
		return nil, diags
	}
	return group, nil
}

func (g *groupExec) Delete(ctx context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse, data *pb.Group) diag.Diagnostics {
	_, err := g.groupClient.DeleteGroup(ctx, data)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("failed to delete group", err.Error())
		return diags
	}
	return nil
}

func (g *groupExec) SetGroupServiceClient(client pb.GroupServiceClient) {
	g.groupClient = client
}

func (g *groupExec) SetUserServiceClient(client pb.UserServiceClient) {
	g.userClient = client
}
