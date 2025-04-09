// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datazone

// **PLEASE DELETE THIS AND ALL TIP COMMENTS BEFORE SUBMITTING A PR FOR REVIEW!**
//
// TIP: ==== INTRODUCTION ====
// Thank you for trying the skaff tool!
//
// You have opted to include these helpful comments. They all include "TIP:"
// to help you find and remove them when you're done with them.
//
// While some aspects of this file are customized to your input, the
// scaffold tool does *not* look at the AWS API and ensure it has correct
// function, structure, and variable names. It makes guesses based on
// commonalities. You will need to make significant adjustments.
//
// In other words, as generated, this is a rough outline of the work you will
// need to do. If something doesn't make sense for your situation, get rid of
// it.

import (
	// TIP: ==== IMPORTS ====
	// This is a common set of imports but not customized to your code since
	// your code hasn't been written yet. Make sure you, your IDE, or
	// goimports -w <file> fixes these imports.
	//
	// The provider linter wants your imports to be in two groups: first,
	// standard library (i.e., "fmt" or "strings"), second, everything else.
	//
	// Also, AWS Go SDK v2 may handle nested structures differently than v1,
	// using the services/datazone/types package. If so, you'll
	// need to import types and reference the nested types, e.g., as
	// awstypes.<Type Name>.
	"context"
	"errors"
    "fmt"
	"time"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go-v2/service/datazone"
	awstypes "github.com/aws/aws-sdk-go-v2/service/datazone/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/fwdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
	"github.com/hashicorp/terraform-provider-aws/internal/framework/flex"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// TIP: ==== FILE STRUCTURE ====
// All resources should follow this basic outline. Improve this resource's
// maintainability by sticking to it.
//
// 1. Package declaration
// 2. Imports
// 3. Main resource struct with schema method
// 4. Create, read, update, delete methods (in that order)
// 5. Other functions (flatteners, expanders, waiters, finders, etc.)

// Function annotations are used for resource registration to the Provider. DO NOT EDIT.
// @FrameworkResource("aws_datazone_project_membership", name="Project Membership")
func newResourceProjectMembership(_ context.Context) (resource.ResourceWithConfigure, error) {
	r := &resourceProjectMembership{}
	
	// TIP: ==== CONFIGURABLE TIMEOUTS ====
	// Users can configure timeout lengths but you need to use the times they
	// provide. Access the timeout they configure (or the defaults) using,
	// e.g., r.CreateTimeout(ctx, plan.Timeouts) (see below). The times here are
	// the defaults if they don't configure timeouts.
	r.SetDefaultCreateTimeout(10 * time.Minute)
	r.SetDefaultDeleteTimeout(10 * time.Minute)

	return r, nil
}

const (
	ResNameProjectMembership = "Project Membership"
)

type resourceProjectMembership struct {
	framework.ResourceWithConfigure
	framework.WithTimeouts
	framework.WithNoUpdate
}


// TIP: ==== SCHEMA ====
// In the schema, add each of the attributes in snake case (e.g.,
// delete_automated_backups).
//
// Formatting rules:
// * Alphabetize attributes to make them easier to find.
// * Do not add a blank line between attributes.
//
// Attribute basics:
// * If a user can provide a value ("configure a value") for an
//   attribute (e.g., instances = 5), we call the attribute an
//   "argument."
// * You change the way users interact with attributes using:
//     - Required
//     - Optional
//     - Computed
// * There are only four valid combinations:
//
// 1. Required only - the user must provide a value
// Required: true,
//
// 2. Optional only - the user can configure or omit a value; do not
//    use Default or DefaultFunc
// Optional: true,
//
// 3. Computed only - the provider can provide a value but the user
//    cannot, i.e., read-only
// Computed: true,
//
// 4. Optional AND Computed - the provider or user can provide a value;
//    use this combination if you are using Default
// Optional: true,
// Computed: true,
//
// You will typically find arguments in the input struct
// (e.g., CreateDBInstanceInput) for the create operation. Sometimes
// they are only in the input struct (e.g., ModifyDBInstanceInput) for
// the modify operation.
//
// For more about schema options, visit
// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/schemas?page=schemas
func (r *resourceProjectMembership) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"domain_identifier": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexache.MustCompile(`^dzd[-_][a-zA-Z0-9_-]{1,36}$`), "must conform to: ^dzd[-_][a-zA-Z0-9_-]{1,36}$ "),
				},
			},
		    "project_identifier": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexache.MustCompile(`^[\w -]+$`), "must conform to: ^[\\w -]+$ "),
					stringvalidator.LengthBetween(1, 64),
				},
			},
            "user_designation": schema.StringAttribute{
                Required: true,
                Validators: []validator.String{
                    stringvalidator.OneOf("PROJECT_OWNER", "PROJECT_CONTRIBUTOR", "PROJECT_VIEWER", "PROJECT_CONSUMER", "PROJECT_STEWARD"),
                },
            },
            "member": schema.StringAttribute{
                Required: true,
            },
		},
		Blocks: map[string]schema.Block{
			names.AttrTimeouts: timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Delete: true,
			}),
		},
	}
}

func (r *resourceProjectMembership) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// TIP: ==== RESOURCE CREATE ====
	// Generally, the Create function should do the following things. Make
	// sure there is a good reason if you don't do one of these.
	//
	// 1. Get a client connection to the relevant service
	// 2. Fetch the plan
	// 3. Populate a create input structure
	// 4. Call the AWS create/put function
	// 5. Using the output from the create function, set the minimum arguments
	//    and attributes for the Read function to work, as well as any computed
	//    only attributes. 
	// 6. Use a waiter to wait for create to complete
	// 7. Save the request plan to response state

	// TIP: -- 1. Get a client connection to the relevant service
	conn := r.Meta().DataZoneClient(ctx)
	
	// TIP: -- 2. Fetch the plan
	var plan resourceProjectMembershipModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TIP: -- 3. Populate a Create input structure
	var input datazone.CreateProjectMembershipInput
	// TIP: Using a field name prefix allows mapping fields such as `ID` to `ProjectMembershipId`
	//resp.Diagnostics.Append(flex.Expand(ctx, plan, &input)...)
	//if resp.Diagnostics.HasError() {
	//	return
	//}

    input.DomainIdentifier = plan.DomainIdentifier.ValueStringPointer()
    input.ProjectIdentifier = plan.ProjectIdentifier.ValueStringPointer()

    input.Designation = awstypes.UserDesignation(plan.UserDesignation.ValueString())

	if !plan.Member.IsNull() {
        input.Member = &awstypes.MemberMemberUserIdentifier{
            Value: plan.Member.ValueString(),
        }
	}

    // Human friendly ID for errors
    id := fmt.Sprintf("%s:%s:%s", *input.DomainIdentifier, *input.ProjectIdentifier, plan.Member.ValueString())

	// TIP: -- 4. Call the AWS Create function
	out, err := conn.CreateProjectMembership(ctx, &input)
	if err != nil {
		// TIP: Since ID has not been set yet, you cannot use plan.ID.String()
		// in error messages at this point.
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.DataZone, create.ErrActionCreating, ResNameProjectMembership, id, err),
			err.Error(),
		)
		return
	}
	if out == nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.DataZone, create.ErrActionCreating, ResNameProjectMembership, id, nil),
			errors.New("failed when creating").Error(),
		)
		return
	}

	// TIP: -- 5. Using the output from the create function, set attributes
	resp.Diagnostics.Append(flex.Flatten(ctx, out, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TIP: -- 6. Use a waiter to wait for create to complete
	createTimeout := r.CreateTimeout(ctx, plan.Timeouts)
    findProjectMembershipInput := &FindProjectMembershipInput{
        DomainIdentifier: plan.DomainIdentifier.ValueStringPointer(),
        ProjectIdentifier: plan.ProjectIdentifier.ValueStringPointer(),
        Member: plan.Member.ValueStringPointer(),
    }
	_, err = waitProjectMembershipCreated(ctx, conn, findProjectMembershipInput, createTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.DataZone, create.ErrActionWaitingForCreation, ResNameProjectMembership, id, err),
			err.Error(),
		)
		return
	}
	
	// TIP: -- 7. Save the request plan to response state
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *resourceProjectMembership) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// TIP: ==== RESOURCE READ ====
	// Generally, the Read function should do the following things. Make
	// sure there is a good reason if you don't do one of these.
	//
	// 1. Get a client connection to the relevant service
	// 2. Fetch the state
	// 3. Get the resource from AWS
	// 4. Remove resource from state if it is not found
	// 5. Set the arguments and attributes
	// 6. Set the state

	// TIP: -- 1. Get a client connection to the relevant service
	conn := r.Meta().DataZoneClient(ctx)
	
	// TIP: -- 2. Fetch the state
	var state resourceProjectMembershipModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	
    findProjectMembershipInput := &FindProjectMembershipInput{
        DomainIdentifier: state.DomainIdentifier.ValueStringPointer(),
        ProjectIdentifier: state.ProjectIdentifier.ValueStringPointer(),
        Member: state.Member.ValueStringPointer(),
    }
	// TIP: -- 3. Get the resource from AWS using an API Get, List, or Describe-
	// type function, or, better yet, using a finder.
	out, err := findProjectMembership(ctx, conn, findProjectMembershipInput)
	// TIP: -- 4. Remove resource from state if it is not found
	if tfresource.NotFound(err) {
		resp.Diagnostics.Append(fwdiag.NewResourceNotFoundWarningDiagnostic(err))
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.DataZone, create.ErrActionReading, ResNameProjectMembership, state.Member.String(), err),
			err.Error(),
		)
		return
	}
	
	// TIP: -- 5. Set the arguments and attributes
	resp.Diagnostics.Append(flex.Flatten(ctx, out, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	
	// TIP: -- 6. Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *resourceProjectMembership) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// TIP: ==== RESOURCE DELETE ====
	// Most resources have Delete functions. There are rare situations
	// where you might not need a delete:
	// a. The AWS API does not provide a way to delete the resource
	// b. The point of your resource is to perform an action (e.g., reboot a
	//    server) and deleting serves no purpose.
	//
	// The Delete function should do the following things. Make sure there
	// is a good reason if you don't do one of these.
	//
	// 1. Get a client connection to the relevant service
	// 2. Fetch the state
	// 3. Populate a delete input structure
	// 4. Call the AWS delete function
	// 5. Use a waiter to wait for delete to complete
	// TIP: -- 1. Get a client connection to the relevant service
	conn := r.Meta().DataZoneClient(ctx)
	
	// TIP: -- 2. Fetch the state
	var state resourceProjectMembershipModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	
	// TIP: -- 3. Populate a delete input structure
	input := datazone.DeleteProjectMembershipInput{
        DomainIdentifier: state.DomainIdentifier.ValueStringPointer(),
        Member: &awstypes.MemberMemberUserIdentifier{
            Value: *state.Member.ValueStringPointer(),
        },
        ProjectIdentifier: state.ProjectIdentifier.ValueStringPointer(),
	}
	
	// TIP: -- 4. Call the AWS delete function
	_, err := conn.DeleteProjectMembership(ctx, &input)
	// TIP: On rare occassions, the API returns a not found error after deleting a
	// resource. If that happens, we don't want it to show up as an error.
	if err != nil {
		if errs.IsA[*awstypes.ResourceNotFoundException](err) {
			return
		}

		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.DataZone, create.ErrActionDeleting, ResNameProjectMembership, state.Member.String(), err),
			err.Error(),
		)
		return
	}
	
	// TIP: -- 5. Use a waiter to wait for delete to complete
	deleteTimeout := r.DeleteTimeout(ctx, state.Timeouts)
    findProjectMembershipInput := &FindProjectMembershipInput{
        DomainIdentifier: state.DomainIdentifier.ValueStringPointer(),
        ProjectIdentifier: state.ProjectIdentifier.ValueStringPointer(),
        Member: state.Member.ValueStringPointer(),
    }
	_, err = waitProjectMembershipDeleted(ctx, conn, findProjectMembershipInput, deleteTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.DataZone, create.ErrActionWaitingForDeletion, ResNameProjectMembership, state.Member.String(), err),
			err.Error(),
		)
		return
	}
}

// TIP: ==== TERRAFORM IMPORTING ====
// If Read can get all the information it needs from the Identifier
// (i.e., path.Root("id")), you can use the PassthroughID importer. Otherwise,
// you'll need a custom import function.
//
// See more:
// https://developer.hashicorp.com/terraform/plugin/framework/resources/import
func (r *resourceProjectMembership) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root(names.AttrID), req, resp)
}


// TIP: ==== STATUS CONSTANTS ====
// Create constants for states and statuses if the service does not
// already have suitable constants. We prefer that you use the constants
// provided in the service if available (e.g., awstypes.StatusInProgress).
const (
	statusChangePending = "Pending"
	statusDeleting      = "Deleting"
	statusNormal        = "Normal"
	statusUpdated       = "Updated"
)

// TIP: ==== WAITERS ====
// Some resources of some services have waiters provided by the AWS API.
// Unless they do not work properly, use them rather than defining new ones
// here.
//
// Sometimes we define the wait, status, and find functions in separate
// files, wait.go, status.go, and find.go. Follow the pattern set out in the
// service and define these where it makes the most sense.
//
// If these functions are used in the _test.go file, they will need to be
// exported (i.e., capitalized).
//
// You will need to adjust the parameters and names to fit the service.
func waitProjectMembershipCreated(ctx context.Context, conn *datazone.Client, findProjectMembershipInput *FindProjectMembershipInput, timeout time.Duration) (*datazone.CreateProjectMembershipOutput, error) {
	stateConf := &retry.StateChangeConf{
		Pending:                   []string{},
		Target:                    []string{statusNormal},
		Refresh:                   statusProjectMembership(ctx, conn, findProjectMembershipInput),
		Timeout:                   timeout,
		NotFoundChecks:            20,
		ContinuousTargetOccurence: 2,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if out, ok := outputRaw.(*datazone.CreateProjectMembershipOutput); ok {
		return out, err
	}

	return nil, err
}

// TIP: A deleted waiter is almost like a backwards created waiter. There may
// be additional pending states, however.
func waitProjectMembershipDeleted(ctx context.Context, conn *datazone.Client, findProjectMembershipInput *FindProjectMembershipInput, timeout time.Duration) (*datazone.DeleteProjectMembershipOutput, error) {
	stateConf := &retry.StateChangeConf{
		Pending:                   []string{statusDeleting, statusNormal},
		Target:                    []string{},
		Refresh:                   statusProjectMembership(ctx, conn, findProjectMembershipInput),
		Timeout:                   timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if out, ok := outputRaw.(*datazone.DeleteProjectMembershipOutput); ok {
		return out, err
	}

	return nil, err
}

// TIP: ==== STATUS ====
// The status function can return an actual status when that field is
// available from the API (e.g., out.Status). Otherwise, you can use custom
// statuses to communicate the states of the resource.
//
// Waiters consume the values returned by status functions. Design status so
// that it can be reused by a create, update, and delete waiter, if possible.
func statusProjectMembership(ctx context.Context, conn *datazone.Client, findProjectMembershipInput *FindProjectMembershipInput) retry.StateRefreshFunc {
	return func() (any, string, error) {
		out, err := findProjectMembership(ctx, conn, findProjectMembershipInput)
		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return out, statusNormal, nil
	}
}

// TIP: ==== FINDERS ====
// The find function is not strictly necessary. You could do the API
// request from the status function. However, we have found that find often
// comes in handy in other places besides the status function. As a result, it
// is good practice to define it separately.
func findProjectMembership(ctx context.Context, conn *datazone.Client, findProjectMembershipInput *FindProjectMembershipInput) (*string, error) {
    input := datazone.ListProjectMembershipsInput{
        DomainIdentifier: findProjectMembershipInput.DomainIdentifier,
        ProjectIdentifier: findProjectMembershipInput.ProjectIdentifier,
    }

	pages := datazone.NewListProjectMembershipsPaginator(conn, &input)

	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)

        if err != nil {
            if errs.IsA[*awstypes.ResourceNotFoundException](err) {
                return nil, &retry.NotFoundError{
                    LastError:   err,
                    LastRequest: &input,
                }
            }

            return nil, err
        }

		for _, membership := range page.Members {
            if group, ok := membership.MemberDetails.(*awstypes.MemberDetailsMemberGroup); ok {
                if findProjectMembershipInput.Member == group.Value.GroupId {
                    return findProjectMembershipInput.Member, nil
                }
            }
            if user, ok := membership.MemberDetails.(*awstypes.MemberDetailsMemberUser); ok {
                if findProjectMembershipInput.Member == user.Value.UserId {
                    return findProjectMembershipInput.Member, nil
                }
            }
		}
	}

    return nil, &retry.NotFoundError{}
}

// TIP: ==== DATA STRUCTURES ====
// With Terraform Plugin-Framework configurations are deserialized into
// Go types, providing type safety without the need for type assertions.
// These structs should match the schema definition exactly, and the `tfsdk`
// tag value should match the attribute name. 
//
// Nested objects are represented in their own data struct. These will 
// also have a corresponding attribute type mapping for use inside flex
// functions.
//
// See more:
// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/accessing-values
type resourceProjectMembershipModel struct {
	DomainIdentifier    types.String                                          `tfsdk:"domain_identifier"`
    UserDesignation     types.String                                          `tfsdk:"user_designation"`
    Member              types.String                                          `tfsdk:"member"`
    ProjectIdentifier   types.String                                          `tfsdk:"project_identifier"`
	Timeouts            timeouts.Value                                        `tfsdk:"timeouts"`
}

type FindProjectMembershipInput struct {
    DomainIdentifier    *string
    ProjectIdentifier   *string
    Member              *string
}

