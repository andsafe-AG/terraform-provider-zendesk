# Initialize a Zendesk provider with the Zendesk API using API-Token Authentication.
provider "zendesk" {
  # host_url is the URL of your Zendesk instance
  host_url = "https://your-zendesk-instance.com/"
  # email is the email address of the Zendesk user you want to use for authentication
  email = "your@email.com"
  # api_token is the API token of the Zendesk user you want to use for authentication
  api_token = "your-api-token"
}