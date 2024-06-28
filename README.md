# terraform-provider-zendesk
[![Tests Status](https://github.com/andsafe-AG/terraform-provider-zendesk/workflows/Tests/badge.svg)](https://github.com/andsafe-AG/terraform-provider-zendesk/actions)


# Zendesk Terraform Provider by andsafe-AG

This Terraform Provider is developed on terraform Plugin Framework and integrates existing Zendesk Terraform provider with [Terraform Provider Zendesk by Nukosuke ](https://registry.terraform.io/providers/nukosuke/zendesk/latest/docs)

28.06.2024 - Combined existing provider with [Terraform Provider Zendesk by Nukosuke (GitHub)](https://github.com/nukosuke/terraform-provider-zendesk)
27.06.2024 - Custom Ticket Status Support


## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= **1.4**
- [Go](https://golang.org/doc/install) >= 1.21


## Using the provider

Configure the provider by setting the `subdomain`, `email` and `token` attributes.

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
