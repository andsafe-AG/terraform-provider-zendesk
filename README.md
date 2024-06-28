# terraform-provider-zendesk
[![Tests Status](https://github.com/andsafe-AG/terraform-provider-zendesk/workflows/Tests/badge.svg)](https://github.com/andsafe-AG/terraform-provider-zendesk/actions)


# Zendesk Terraform Provider by andsafe-AG

27.06.2024 - Custom Ticket Status Support


## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= **1.4**
- [Go](https://golang.org/doc/install) >= 1.21


## Using the provider

Configure the provider by setting the `subdomain`, `email` and `api_token` attributes.

```hcl
terraform {
  required_providers {
    zendesk = {
      source  = "andsafe-AG/zendesk"
      version = ">= 0.0.1"
    }
  }
}
```

## License
MPL 2.0 License
