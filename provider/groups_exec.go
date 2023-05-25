package main

import (
	"context"
	"github.com/travix/gotf-example/provider/providerpb"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/travix/gotf-example/pb"
)

var _ providerpb.GroupsDataSourceExec = &groupsExec{}

type groupsExec struct {
	client pb.GroupServiceClient
}

func (g *groupsExec) SetGroupServiceClient(client pb.GroupServiceClient) {
	g.client = client
}

func (g *groupsExec) Read(ctx context.Context, _ datasource.ReadRequest, _ *datasource.ReadResponse, _ *pb.Groups) (*pb.Groups, diag.Diagnostics) {
	groups, err := g.client.ListGroups(ctx, &pb.Empty{})
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("failed to list groups", err.Error())
		return nil, diags
	}
	for _, group := range groups.Groups {
		for _, user := range group.Users {
			group.UsersNames = append(group.UsersNames, user.Username)
		}
	}
	return groups, nil
}
