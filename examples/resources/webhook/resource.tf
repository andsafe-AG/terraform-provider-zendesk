# Webhook resource
# For API Details see https://developer.zendesk.com/api-reference/webhooks/webhooks-api/webhooks/

resource "zendesk_webhook" "my_webhook" {
  webhook = {
    name = "My Webhook"
    # Authentication is optional
    authentication = {
      # basic_auth or api_key
      type = "basic_auth"
      # required
      add_position = "header"
      data = {
        # Use username and password for Basic Authentication.
        username = "my-username"
        password = "my-password"
        # Alternatively, use token attribute for API Key Authentication.
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
    # External source by which a webhook is created, e.g. Zendesk Marketplace. Optional
    external_source = {
      external_source_data = {
        app_id          = "my-app-id"
        installation_id = "my-installation-id"
      },
      type = "zendesk_app"
    }
    # Required. Current status of the webhook. Allowed values are "active", or "inactive".
    status = "active"
    # Event subscriptions for the webhook.
    # To subscribe the webhook to Zendesk events, specify one or more event types.
    # For supported event type values, see Webhook event types.
    # To connect the webhook to a trigger or automation, specify only "conditional_ticket_events" in the array.
    subscriptions = [
      "conditional_ticket_events"
    ]

  }
}