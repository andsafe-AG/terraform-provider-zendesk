---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "zendesk Provider"
subcategory: ""
description: |-
  The Zendesk provider allows you to interact with the Zendesk API.
---

# zendesk Provider

The Zendesk provider allows you to interact with the Zendesk API.

## Example Usage

```terraform
# Initialize a Zendesk provider with the Zendesk API using API-Token Authentication.
provider "zendesk" {
  # host_url is the URL of your Zendesk instance
  host_url = "https://your-zendesk-instance.com/"
  # email is the email address of the Zendesk user you want to use for authentication
  email = "your@email.com"
  # api_token is the API token of the Zendesk user you want to use for authentication
  api_token = "your-api-token"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `api_token` (String, Sensitive) The API token to authenticate with.
- `email` (String, Sensitive) The email address of the user to authenticate with. It will be masked.
- `host_url` (String) The base URL of your Zendesk instance.
