# Initialize a Zendesk provider with the Zendesk API using API-Token Authentication.
provider "zendesk" {
  # subdomain of your Zendesk instance in https://subdomain.zendesk.com/
  account = "your-zendesk-instance"
  # email is the email address of the Zendesk user you want to use for authentication
  email = "your@email.com"
  # token is the API token of the Zendesk user you want to use for authentication
  token = "your-api-token"
}