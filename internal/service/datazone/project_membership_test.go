// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datazone_test
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
	// types.<Type Name>.
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/datazone"
	"github.com/aws/aws-sdk-go-v2/service/datazone/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/names"

	// TIP: You will often need to import the package that this test file lives
	// in. Since it is in the "test" context, it must import the package to use
	// any normal context constants, variables, or functions.
	tfdatazone "github.com/hashicorp/terraform-provider-aws/internal/service/datazone"
)

// TIP: File Structure. The basic outline for all test files should be as
// follows. Improve this resource's maintainability by following this
// outline.
//
// 1. Package declaration (add "_test" since this is a test file)
// 2. Imports
// 3. Unit tests
// 4. Basic test
// 5. Disappears test
// 6. All the other tests
// 7. Helper functions (exists, destroy, check, etc.)
// 8. Functions that return Terraform configurations


// TIP: ==== ACCEPTANCE TESTS ====
// This is an example of a basic acceptance test. This should test as much of
// standard functionality of the resource as possible, and test importing, if
// applicable. We prefix its name with "TestAcc", the service, and the
// resource name.
//
// Acceptance test access AWS and cost money to run.
func TestAccDataZoneProjectMembership_basic(t *testing.T) {
	ctx := acctest.Context(t)
	// TIP: This is a long-running test guard for tests that run longer than
	// 300s (5 min) generally.
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	pName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
    rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	projectName := "aws_datazone_project.test"
	domainName := "aws_datazone_domain.test"
	resourceName := "aws_datazone_project_membership.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.DataZoneEndpointID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.DataZoneServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProjectMembershipDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccProjectMembershipConfig_basic(pName, dName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProjectMembershipExists(ctx, dName, pName, rName),
					resource.TestCheckResourceAttr(resourceName, "auto_minor_version_upgrade", "false"),
					resource.TestCheckResourceAttr(resourceName, "domain_identifier", domainName),
					resource.TestCheckResourceAttr(resourceName, "project_identifier", projectName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"apply_immediately", "user"},
			},
		},
	})
}

func TestAccDataZoneProjectMembership_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var projectmembership datazone.DescribeProjectMembershipResponse
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_datazone_project_membership.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.DataZoneEndpointID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.DataZoneServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProjectMembershipDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccProjectMembershipConfig_basic(rName, testAccProjectMembershipVersionNewer),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProjectMembershipExists(ctx, resourceName, &projectmembership),
					// TIP: The Plugin-Framework disappears helper is similar to the Plugin-SDK version,
					// but expects a new resource factory function as the third argument. To expose this
					// private function to the testing package, you may need to add a line like the following
					// to exports_test.go:
					//
					//   var ResourceProjectMembership = newResourceProjectMembership
					acctest.CheckFrameworkResourceDisappears(ctx, acctest.Provider, tfdatazone.ResourceProjectMembership, resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckProjectMembershipDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).DataZoneClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_datazone_project_membership" {
				continue
			}

			// TIP: ==== FINDERS ====
			// The find function should be exported. Since it won't be used outside of the package, it can be exported
			// in the `exports_test.go` file.
			_, err := tfdatazone.FindProjectMembershipByID(ctx, conn, rs.Primary.ID)
			if tfresource.NotFound(err) {
				return nil
			}
			if err != nil {
			        return create.Error(names.DataZone, create.ErrActionCheckingDestroyed, tfdatazone.ResNameProjectMembership, rs.Primary.ID, err)
			}

			return create.Error(names.DataZone, create.ErrActionCheckingDestroyed, tfdatazone.ResNameProjectMembership, rs.Primary.ID, errors.New("not destroyed"))
		}

		return nil
	}
}

func testAccCheckProjectMembershipExists(ctx context.Context, domainIdentifier, projectIdentifier, memberName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[memberName]
		if !ok {
			return create.Error(names.DataZone, create.ErrActionCheckingExistence, tfdatazone.ResNameProjectMembership, memberName, errors.New("not found"))
		}

		if rs.Primary.ID == "" {
			return create.Error(names.DataZone, create.ErrActionCheckingExistence, tfdatazone.ResNameProjectMembership, memberName, errors.New("not set"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).DataZoneClient(ctx)

        findProjectMembershipInput := &FindProjectMembershipInput{
            DomainIdentifier: domainIdentifier,
            ProjectIdentifier: projectIdentifier,
            Member: memberName,
        }

        _, err := tfdatazone.FindProjectMembership(ctx, conn, findProjectMembershipInput)

		if err != nil {
			return create.Error(names.DataZone, create.ErrActionCheckingExistence, tfdatazone.ResNameProjectMembership, rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccCheckProjectMembershipNotRecreated(before, after *datazone.DescribeProjectMembershipResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before, after := aws.ToString(before.ProjectMembershipId), aws.ToString(after.ProjectMembershipId); before != after {
			return create.Error(names.DataZone, create.ErrActionCheckingNotRecreated, tfdatazone.ResNameProjectMembership, aws.ToString(before.ProjectMembershipId), errors.New("recreated"))
		}

		return nil
	}
}

func testAccProjectMembershipConfig_basic(pName, dName, roleName string) string {
    return acctest.ConfigCompose(testAccProjectConfig_basic(pName, dName), fmt.Sprintf(`
resource "aws_iam_role" "test" {
  name = %[1]q
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "ec2.amazonaws.com"
        ]
      },
      "Action": [
        "sts:AssumeRole"
      ]
    }
  ]
}
EOF
}

resource "aws_datazone_project_membership" "test" {
  domain_identifier   = aws_datazone_domain.test.id
  project_identifier  = aws_datazone_project.test.id
  member              = aws_iam_role.test.arn
}
`, roleName))
}

