// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connect

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/connect"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @FrameworkResource(name="queue Quick Connect Association")
func newResourcequeueQuickConnectAssociation(_ context.Context) (resource.ResourceWithConfigure, error) {
	r := &resourcequeueQuickConnectAssociation{}
	r.SetDefaultCreateTimeout(30 * time.Minute)
	r.SetDefaultUpdateTimeout(30 * time.Minute)
	r.SetDefaultDeleteTimeout(30 * time.Minute)

	return r, nil
}

const (
	ResNamequeueQuickConnectAssociation = "queue Quick Connect Association"
)

type resourceQueueQuickConnectAssociationData struct {
	ID              types.String `tfsdk:"id"`
	InstanceId      types.String `tfsdk:"instance_id"`
	QueueId         types.String `tfsdk:"queue_id"`
	QuickConnectIds []*string    `tfsdk:"quick_connect_ids"`
}

type resourcequeueQuickConnectAssociation struct {
	framework.ResourceWithConfigure
	framework.WithTimeouts
}

func (r *resourcequeueQuickConnectAssociation) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "aws_connect_queue_quick_connect"
}

func (r *resourcequeueQuickConnectAssociation) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttribute(),
			"instance_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"queue_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"quick_connect_ids": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *resourcequeueQuickConnectAssociation) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceQueueQuickConnectAssociationData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	conn := r.Meta().ConnectConn(ctx)
	in := &connect.AssociateQueueQuickConnectsInput{
		InstanceId:      aws.String(plan.InstanceId.ValueString()),
		QueueId:         aws.String(plan.QueueId.ValueString()),
		QuickConnectIds: plan.QuickConnectIds,
	}

	_, err := conn.AssociateQueueQuickConnectsWithContext(ctx, in)
	if err != nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(
				names.Connect,
				create.ErrActionCreating,
				ResNamequeueQuickConnectAssociation,
				plan.InstanceId.String(),
				err,
			),
			err.Error(),
		)
	}

	plan.ID = types.StringValue(fmt.Sprintf("%v:%v", plan.InstanceId.ValueString(), plan.QueueId.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *resourcequeueQuickConnectAssociation) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceQueueQuickConnectAssociationData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	conn := r.Meta().ConnectConn(ctx)
	in := &connect.ListQueueQuickConnectsInput{
		InstanceId: aws.String(state.InstanceId.ValueString()),
		QueueId:    aws.String(state.QueueId.ValueString()),
	}
	quickConnectIds, err := ListAllQueueQuickConnectIdsWithContext(ctx, conn, in)
	if err != nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(
				names.Connect,
				create.ErrActionReading,
				ResNamequeueQuickConnectAssociation,
				state.InstanceId.String(),
				err,
			),
			err.Error(),
		)
		return
	}

	state.QuickConnectIds = append(state.QuickConnectIds, quickConnectIds...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func ListAllQueueQuickConnectIdsWithContext(ctx context.Context, conn *connect.Connect, in *connect.ListQueueQuickConnectsInput) ([]*string, error) {
	quickConnects := []*string{}
	for {
		response, err := conn.ListQueueQuickConnectsWithContext(ctx, in)
		if err != nil {
			return nil, err
		}
		for _, quickConnect := range response.QuickConnectSummaryList {
			quickConnects = append(quickConnects, quickConnect.Id)
		}
		if response.NextToken == nil {
			break
		}
	}
	return quickConnects, nil
}

func (r *resourcequeueQuickConnectAssociation) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *resourcequeueQuickConnectAssociation) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

// TIP: ==== TERRAFORM IMPORTING ====
// If Read can get all the information it needs from the Identifier
// (i.e., path.Root("id")), you can use the PassthroughID importer. Otherwise,
// you'll need a custom import function.
//
// See more:
// https://developer.hashicorp.com/terraform/plugin/framework/resources/import
func (r *resourcequeueQuickConnectAssociation) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
