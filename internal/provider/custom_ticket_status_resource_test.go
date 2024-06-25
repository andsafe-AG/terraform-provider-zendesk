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
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomStatusResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + testAccCustomStatusResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zendesk_custom_status.test", "custom_status.agent_label", "one"),
					resource.TestCheckResourceAttr("zendesk_custom_status.test", "custom_status.active", "true"),
					resource.TestCheckResourceAttrSet("zendesk_custom_status.test", "custom_status_id"),
				),
			},
			// ImportState testing
			{
				ResourceName:                         "zendesk_custom_status.test",
				ImportState:                          true,
				ImportStateVerifyIdentifierAttribute: "custom_status_id",
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"agent_label"},
			},
			// Update and Read testing
			{
				Config: providerConfig + testAccCustomStatusResourceConfig("two"),

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"zendesk_custom_status.test", "custom_status.agent_label", "two"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCustomStatusResourceConfig(agentLabel string) string {
	return fmt.Sprintf(`
resource "zendesk_custom_status" "test" {
custom_status = {
    status_category: "open"
  	agent_label: %[1]q
	end_user_label: %[1]q
    active: true
  }
}
`, agentLabel)
}
