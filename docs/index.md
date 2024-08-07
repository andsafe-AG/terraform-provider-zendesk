---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "zendesk Provider"
subcategory: ""
description: |-
  
---

# zendesk Provider



## Example Usage

```terraform
# Initialize a Zendesk provider with the Zendesk API using API-Token Authentication.
provider "zendesk" {
  # subdomain of your Zendesk instance in https://subdomain.zendesk.com/
  account = "your-zendesk-instance"
  # email is the email address of the Zendesk user you want to use for authentication
  email = "your@email.com"
  # token is the API token of the Zendesk user you want to use for authentication
  token = "your-api-token"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `account` (String) Account name of your Zendesk instance.
- `email` (String) Email address of agent user who have permission to access the API.
- `token` (String, Sensitive) [API token](https://developer.zendesk.com/rest_api/docs/support/introduction#api-token) for your Zendesk instance.
