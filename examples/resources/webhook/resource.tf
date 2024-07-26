# Webhook resource
# For API Details see https://developer.zendesk.com/api-reference/webhooks/webhooks-api/webhooks/
# Signing secret will be fetched to the state after creation, it cannot be set.
resource "zendesk_webhook" "my_webhook" {
  webhook = {
    name = "My Webhook"

    # Authentication is optional
    authentication = {
      # Currently supported only Basic Authentication and Bearer Token.
      # The Zendesk API supports API Key, though not documented in the API Reference.
      # basic_auth or bearer_token
      type = "basic_auth"
      # required
      add_position = "header"
      data = {
        # username and password are required for Basic Authentication.
        username = "my-username"
        password = "my-password"
        # Alternatively, use 'token' attribute for Bearer Token Authentication.
      }
    }

    # The destination URL that the webhook notifies when Zendesk events occur
    endpoint = "https://example.com/webhook"

    # Required. Allowed values are "GET", "POST", "PUT", "PATCH", or "DELETE"
    http_method = "POST"

    # Required. Allowed values are "json", "xml", or "form_encoded".
    request_format = "json"

    # Optional
    custom_headers = {
      "X-My-Header" = "My-Value"
    }
    description = "My Webhook Description"

    # Required. Current status of the webhook. Allowed values are "active", or "inactive".
    status = "active"
    # Event subscriptions for the webhook.
    # To subscribe the webhook to Zendesk events, specify one or more event types.
    # For supported event type values, see Webhook event types.
    # To connect the webhook to a trigger or automation, specify only "conditional_ticket_events" in the array.
    subscriptions = [
      "conditional_ticket_events"
    ]

    # Attributes external_source and signing_secret cannot be set,
    # but might be returned by the API after the webhook is created.
  }
}