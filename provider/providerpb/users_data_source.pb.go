// Code generated by protoc-gen-terraform. DO NOT EDIT.
// versions:
//   protoc-gen-gotf 0.4.0

package providerpb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/travix/gotf"
	"github.com/travix/gotf/dtsrc"
	"google.golang.org/grpc"

	pb "github.com/travix/gotf-example/pb"
)

// Ensure *UsersDataSource fully satisfy terraform framework interfaces.
var _ datasource.DataSource = &UsersDataSource{}

type UsersDataSourceExec interface {
	dtsrc.Datasource[*pb.Users]
	SetUserServiceClient(pb.UserServiceClient)
}

type UsersDataSource struct {
	exec UsersDataSourceExec
}

func NewUsersDataSource(exec UsersDataSourceExec) func() datasource.DataSource {
	if exec == nil {
		panic("UsersDataSourceExec is required")
	}
	return func() datasource.DataSource {
		return &UsersDataSource{exec: exec}
	}
}

func (d *UsersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
	if _exec, ok := d.exec.(dtsrc.CanMetadata); ok {
		_exec.Metadata(ctx, req, resp)
	}
}

func (d *UsersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "terraform datasource",
		Attributes:          (&pb.Users{}).DatasourceSchema(),
	}
	if _exec, ok := d.exec.(dtsrc.CanSchema); ok {
		_exec.Schema(ctx, req, resp)
	}
}

func (d *UsersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Get the gRPC client connection from the ProviderData
	if req.ProviderData != nil {
		conn, ok := req.ProviderData.(grpc.ClientConnInterface)
		if !ok {
			resp.Diagnostics.AddError(
				"Unexpected ProviderData Type",
				fmt.Sprintf("Expected grpc.ClientConnInterface, got: %T. Please report this issue to the provider developers.", req.ProviderData),
			)
			return
		}
		// set the service clients
		d.exec.SetUserServiceClient(pb.NewUserServiceClient(conn))
	}
	if _exec, ok := d.exec.(dtsrc.CanConfigure); ok {
		_exec.Configure(ctx, req, resp)
		return
	}
}

func (d *UsersDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	if _exec, ok := d.exec.(dtsrc.CanConfigValidators); ok {
		return _exec.ConfigValidators(ctx)
	}
	tflog.Warn(ctx, "ConfigValidators method not implemented. Make sure argument to NewUsersDataSource() implements dtsrc.CanConfigValidators interface")
	return nil
}

func (d *UsersDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	if _exec, ok := d.exec.(dtsrc.CanValidateConfig); ok {
		_exec.ValidateConfig(ctx, req, resp)
		return
	}
	tflog.Warn(ctx, "ValidateConfig method not implemented. Make sure argument to NewUsersDataSource() implements dtsrc.CanValidateConfig interface")
}

func (d *UsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data, diagnostics := gotf.GetModel[pb.Users](ctx, req.Config.Raw, req.Config.Get)
	if diagnostics.HasError() {
		resp.Diagnostics.Append(diagnostics...)
		return
	}
	read, diagnostics := d.exec.Read(ctx, req, resp, data)
	if diagnostics.HasError() {
		resp.Diagnostics.Append(diagnostics...)
		return
	}
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, read)...)
}
