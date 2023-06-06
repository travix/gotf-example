package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	pb "github.com/travix/gotf-example/pb"
	providerpb "github.com/travix/gotf-example/provider/providerpb"
)

// This file was generated by protoc-gen-gotf as a scaffold, it can be modified.
// If you want to regenerate the scaffold delete or rename this file and run protoc with protoc-gen-gotf again.

var _ providerpb.GroupsDataSourceExec = &GroupsDataSourceExec{}

type GroupsDataSourceExec struct {
	groupClient pb.GroupServiceClient
}

func (e *GroupsDataSourceExec) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse, data *pb.Groups) (*pb.Groups, diag.Diagnostics) {
	groups, err := e.groupClient.ListGroups(ctx, &pb.Empty{})
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("failed to get groups", err.Error())
		return nil, diags
	}
	return groups, nil
}

func (e *GroupsDataSourceExec) SetGroupServiceClient(client pb.GroupServiceClient) {
	e.groupClient = client
}
