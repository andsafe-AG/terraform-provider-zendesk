// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

/**
 * To run this test, you need to run the mock server first:
 * 'cd zendesk_api/mock'
 * 'docker build -t zendesk-mock .'
 * 'docker run -it -p 8080:8080 zendesk-mock'
 */

package provider

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomRoleResource(t *testing.T) {

	resource.Test(t, resource.TestCase{

		PreCheck: func() {
			mockContainer, err := startMockServer()
			if err != nil {
				log.Default().Printf("Failed to start container %s\n", err)
				t.Fatal(err)
			}
			t.Cleanup(func() {
				log.Default().Printf("Stopping container %s\n", mockContainer.URI)
				if err := mockContainer.Terminate(context.Background()); err != nil {
					t.Fatalf("failed to terminate container: %s", err)
				}
			})
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read testing
			{
				Config: providerConfig + testAccCustomRoleResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zendesk_custom_role.test", "custom_role_id", "10127"),
					resource.TestCheckResourceAttr("zendesk_custom_role.test", "custom_role.id", "10127"),
					resource.TestCheckResourceAttr("zendesk_custom_role.test", "custom_role.name", "one"),
					resource.TestCheckResourceAttr("zendesk_custom_role.test", "custom_role.description", "test role descr"),
					resource.TestCheckResourceAttr("zendesk_custom_role.test", "custom_role.role_type", "0"),
					resource.TestCheckResourceAttr("zendesk_custom_role.test", "custom_role.configuration.ticket_access", "within-groups"),
					resource.TestCheckResourceAttr("zendesk_custom_role.test", "custom_role.configuration.ticket_editing", "true"),
					resource.TestCheckResourceAttr("zendesk_custom_role.test", "custom_role.configuration.chat_access", "true"),
					resource.TestCheckResourceAttr("zendesk_custom_role.test", "custom_role.configuration.manage_triggers", "false"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + testAccCustomRoleResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zendesk_custom_role.test", "custom_role_id", "10127"),
					resource.TestCheckResourceAttr("zendesk_custom_role.test", "custom_role.name", "two"),
				),
			},
			// ImportState testing by name
			{
				ResourceName:                         "zendesk_custom_role.test",
				ImportState:                          true,
				ImportStateId:                        "two",
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "custom_role_id",
			},
			// ImportState testing by id
			{
				ResourceName:                         "zendesk_custom_role.test",
				ImportState:                          true,
				ImportStateId:                        "10127",
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "custom_role_id",
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCustomRoleResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "zendesk_custom_role" "test" {
  custom_role = {
    name        = %[1]q
    description = "test role descr"
    configuration = {
      ticket_access  = "within-groups"
      ticket_editing = true
    }
  }
}
`, name)
}
