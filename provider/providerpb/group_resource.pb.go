// Code generated by protoc-gen-terraform. DO NOT EDIT.
// versions:
//   protoc-gen-gotf 0.4.0

package providerpb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/travix/gotf"
	"github.com/travix/gotf/rsrc"
	"google.golang.org/grpc"

	pb "github.com/travix/gotf-example/pb"
)

// Ensure *GroupResourceResource fully satisfy terraform framework interfaces.
var _ resource.Resource = &GroupResourceResource{}

type GroupResourceExec interface {
	rsrc.Resource[*pb.Group]
	SetGroupServiceClient(pb.GroupServiceClient)
	SetUserServiceClient(pb.UserServiceClient)
}

type GroupResourceResource struct {
	exec GroupResourceExec
}

func NewGroupResource(exec GroupResourceExec) func() resource.Resource {
	if exec == nil {
		panic("GroupResourceExec is required")
	}
	return func() resource.Resource {
		return &GroupResourceResource{exec: exec}
	}
}

func (r *GroupResourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
	if _exec, ok := r.exec.(rsrc.CanMetadata); ok {
		_exec.Metadata(ctx, req, resp)
	}
}

func (r *GroupResourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "terraform resource",
		Attributes:          (&pb.Group{}).ResourceSchema(),
	}
	if _exec, ok := r.exec.(rsrc.CanSchema); ok {
		_exec.Schema(ctx, req, resp)
	}
}

func (r *GroupResourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
		// Create the service clients
		r.exec.SetGroupServiceClient(pb.NewGroupServiceClient(conn))
		r.exec.SetUserServiceClient(pb.NewUserServiceClient(conn))
	}
	if _exec, ok := r.exec.(rsrc.CanConfigure); ok {
		_exec.Configure(ctx, req, resp)
		return
	}
}

func (r *GroupResourceResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	if _exec, ok := r.exec.(rsrc.CanConfigValidators); ok {
		return _exec.ConfigValidators(ctx)
	}
	tflog.Warn(ctx, "ConfigValidators method not implemented.Make sure argument to NewGroupResourceResource() implements rsrc.CanConfigValidators interface")
	return nil
}

func (r *GroupResourceResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	if _exec, ok := r.exec.(rsrc.CanValidateConfig); ok {
		_exec.ValidateConfig(ctx, req, resp)
		return
	}
	tflog.Warn(ctx, "ValidateConfig method not implemented.Make sure argument to NewGroupResourceResource() implements rsrc.CanValidateConfig interface")
}

func (r *GroupResourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if _exec, ok := r.exec.(rsrc.CanImportState); ok {
		_exec.ImportState(ctx, req, resp)
		return
	}
	tflog.Warn(ctx, "ImportState method not implemented.Make sure argument to NewGroupResourceResource() implements rsrc.CanImportState interface")
}

func (r *GroupResourceResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if _exec, ok := r.exec.(rsrc.CanModifyPlan); ok {
		_exec.ModifyPlan(ctx, req, resp)
		return
	}
	tflog.Warn(ctx, "ModifyPlan method not implemented.Make sure argument to NewGroupResourceResource() implements rsrc.CanModifyPlan interface")
}

func (r *GroupResourceResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	if _exec, ok := r.exec.(rsrc.CanUpgradeState); ok {
		return _exec.UpgradeState(ctx)
	}
	tflog.Warn(ctx, "UpgradeState method not implemented.Make sure argument to NewGroupResourceResource() implements rsrc.CanUpgradeState interface")
	return nil
}

func (r *GroupResourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data, diagnostics := gotf.GetModel[pb.Group](ctx, req.State.Raw, req.State.Get)
	if diagnostics.HasError() {
		resp.Diagnostics.Append(diagnostics...)
		return
	}
	read, diagnostics := r.exec.Read(ctx, req, resp, data)
	if diagnostics.HasError() {
		resp.Diagnostics.Append(diagnostics...)
		return
	}
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, read)...)
}

func (r *GroupResourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data, diagnostics := gotf.GetModel[pb.Group](ctx, req.Config.Raw, req.Config.Get)
	if diagnostics.HasError() {
		resp.Diagnostics.Append(diagnostics...)
		return
	}
	created, diagnostics := r.exec.Create(ctx, req, resp, data)
	if diagnostics.HasError() {
		resp.Diagnostics.Append(diagnostics...)
		return
	}
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, created)...)
}

func (r *GroupResourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data, diagnostics := gotf.GetModel[pb.Group](ctx, req.Config.Raw, req.Config.Get)
	if diagnostics.HasError() {
		resp.Diagnostics.Append(diagnostics...)
		return
	}
	updated, diagnostics := r.exec.Update(ctx, req, resp, data)
	if diagnostics.HasError() {
		resp.Diagnostics.Append(diagnostics...)
		return
	}
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, updated)...)
}

func (r *GroupResourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data, diagnostics := gotf.GetModel[pb.Group](ctx, req.State.Raw, req.State.Get)
	if diagnostics.HasError() {
		resp.Diagnostics.Append(diagnostics...)
		return
	}
	resp.Diagnostics.Append(r.exec.Delete(ctx, req, resp, data)...)
}
