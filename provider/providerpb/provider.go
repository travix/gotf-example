// Code generated by protoc-gen-terraform. DO NOT EDIT.
// versions:
//   protoc-gen-gotf 0.1.3

package providerpb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/travix/gotf/prvdr"

	pb "github.com/travix/gotf-example/pb"
)

// Ensure *ExampleProvider fully satisfy terraform framework interfaces.
var _ provider.Provider = &ExampleProvider{}

type ExampleExec interface {
	prvdr.Provider
	prvdr.CanConfigureGrpc[*pb.ProviderModel]
}

type ExampleProvider struct {
	version string
	exec    ExampleExec
}

func New(version string, exec ExampleExec) func() provider.Provider {
	if exec == nil {
		panic("ExampleExec is required")
	}
	return func() provider.Provider {
		return &ExampleProvider{
			version: version,
			exec:    exec,
		}
	}
}

func (p *ExampleProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "example"
	resp.Version = p.version
	if _exec, ok := p.exec.(prvdr.CanMetadata); ok {
		_exec.Metadata(ctx, req, resp)
	}
}

func (p *ExampleProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "",
		Attributes:  (&pb.ProviderModel{}).ProviderSchema(),
	}
	if _exec, ok := p.exec.(prvdr.CanSchema); ok {
		_exec.Schema(ctx, req, resp)
	}
}

func (p *ExampleProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	if _exec, ok := p.exec.(prvdr.CanConfigure); ok {
		_exec.Configure(ctx, req, resp)
		if resp.DataSourceData != nil {
			resp.Diagnostics.AddWarning("resp.DataSourceData not set", "DataSourceData should be set to grpc.ClientConnInterface by Configure method found nil")
		}
		if resp.ResourceData != nil {
			resp.Diagnostics.AddWarning("resp.ResourceData not set", "ResourceData should be set to grpc.ClientConnInterface by Configure method found nil")
		}
		return
	}
	var data pb.ProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	conn, diagnostics := p.exec.ConfigureGrpc(ctx, &data)
	if diagnostics.HasError() {
		resp.Diagnostics.Append(diagnostics...)
		return
	}
	resp.DataSourceData = conn
	resp.ResourceData = conn
	return

}

func (p *ExampleProvider) ConfigValidators(ctx context.Context) []provider.ConfigValidator {
	if _exec, ok := p.exec.(prvdr.CanConfigValidators); ok {
		return _exec.ConfigValidators(ctx)
	}
	tflog.Warn(ctx, "ConfigValidators method not implemented. Make sure argument to New() implements prvdr.CanConfigValidators interface")
	return nil
}

func (p *ExampleProvider) MetaSchema(ctx context.Context, req provider.MetaSchemaRequest, resp *provider.MetaSchemaResponse) {
	if _exec, ok := p.exec.(prvdr.CanMetaSchema); ok {
		_exec.MetaSchema(ctx, req, resp)
		return
	}
	tflog.Warn(ctx, "MetaSchema method not implemented. Make sure argument to New() implements prvdr.CanMetaSchema interface")
}

func (p *ExampleProvider) ValidateConfig(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) {
	if _exec, ok := p.exec.(prvdr.CanValidateConfig); ok {
		_exec.ValidateConfig(ctx, req, resp)
		return
	}
	tflog.Warn(ctx, "ValidateConfig method not implemented. Make sure argument to New() implements prvdr.CanValidateConfig interface")
}

func (p *ExampleProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return p.exec.DataSources(ctx)
}

func (p *ExampleProvider) Resources(ctx context.Context) []func() resource.Resource {
	return p.exec.Resources(ctx)
}