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
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/testcontainers/testcontainers-go"
)

func TestAccCustomStatusResource(t *testing.T) {

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

type zendeskAPIMockContainer struct {
	testcontainers.Container
	URI string
}

func startMockServer() (*zendeskAPIMockContainer, error) {
	ctx := context.Background()
	dl := log.Default()
	mockContainerName := "zendesk-mock"

	path, err := getDockerfilePath()
	if err != nil {
		return nil, err
	}

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:       path,
			Dockerfile:    "Dockerfile",
			PrintBuildLog: true,
		},
		Name:         mockContainerName,
		ExposedPorts: []string{"8080:8080/tcp"},
		LifecycleHooks: []testcontainers.ContainerLifecycleHooks{
			testcontainers.DefaultLoggingHook(dl),
		},
		WaitingFor: wait.ForLog("Mock engine up and running").WithStartupTimeout(90 * time.Second), // Container is ready
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}
	log.Default().Printf(" Host is: %v \n", ip)

	mappedPort, err := container.MappedPort(ctx, "8080")
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())

	return &zendeskAPIMockContainer{Container: container, URI: uri}, nil
}

func getDockerfilePath() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return "", err
	}
	const ProviderPath = "internal/provider"
	const MockDockerfilePath = "zendesk_api/mock"
	path = strings.Replace(path, ProviderPath, MockDockerfilePath, 1)
	return path, err
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
