package project_ssh_key

import (
	"context"
	"fmt"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name:   "equinix_metal_project_ssh_key",
				Schema: GetResourceSchema(),
			},
		),
	}
}

type Resource struct {
	framework.BaseResource
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {

	r.Meta.AddFwModuleToMetalGoUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metalgo

	// Retrieve values from plan
	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	createRequest := &metalv1.SSHKeyCreateInput{
		Label: plan.Name.ValueStringPointer(),
		Key:   plan.PublicKey.ValueStringPointer(),
	}

	projectId := plan.ProjectID.ValueString()

	// Create API resource
	key, _, err := client.SSHKeysApi.CreateProjectSSHKey(context.Background(), projectId).SSHKeyCreateInput(*createRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create Project SSH Key",
			equinix_errors.FriendlyError(err).Error(),
		)
		return
	}

	// Parse API response into the Terraform state
	resp.Diagnostics.Append(plan.parse(key)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	r.Meta.AddFwModuleToMetalGoUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metalgo

	// Retrieve values from state
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract the ID of the resource from the state
	id := state.ID.ValueString()

	// Use API client to get the current state of the resource
	key, _, err := client.SSHKeysApi.FindSSHKeyById(context.Background(), id).Include(nil).Execute()
	if err != nil {
		err = equinix_errors.FriendlyError(err)

		// If the key is somehow already destroyed, mark as
		// succesfully gone
		if equinix_errors.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"Equinix Metal Project SSHKey not found during refresh",
				fmt.Sprintf("[WARN] SSHKey (%s) not found, removing from state", id),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get Project SSHKey %s", id),
			err.Error(),
		)
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(state.parse(key)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	r.Meta.AddFwModuleToMetalGoUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metalgo

	// Retrieve values from plan
	var state, plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract the ID of the resource from the state
	id := plan.ID.ValueString()

	updateRequest := &metalv1.SSHKeyInput{}
	if !state.Name.Equal(plan.Name) {
		updateRequest.Label = plan.Name.ValueStringPointer()
	}
	if !state.PublicKey.Equal(plan.PublicKey) {
		updateRequest.Key = plan.PublicKey.ValueStringPointer()
	}

	// Update the resource
	key, _, err := client.SSHKeysApi.UpdateSSHKey(context.Background(), id).SSHKeyInput(*updateRequest).Execute()
	if err != nil {
		err = equinix_errors.FriendlyError(err)
		resp.Diagnostics.AddError(
			"Error updating resource",
			"Could not update resource with ID "+id+": "+err.Error(),
		)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(plan.parse(key)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the updated state back into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	r.Meta.AddFwModuleToMetalGoUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metalgo

	// Retrieve values from plan
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract the ID of the resource from the state
	id := state.ID.ValueString()

	// Use API client to delete the resource
	deleteResp, err := client.SSHKeysApi.DeleteSSHKey(context.Background(), id).Execute()
	if equinix_errors.IgnoreHttpResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(deleteResp, err) != nil {
		err = equinix_errors.FriendlyError(err)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete Project SSHKey %s", id),
			err.Error(),
		)
	}
}
