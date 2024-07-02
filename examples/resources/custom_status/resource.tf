# Custom Ticket Status to be used in your Zendesk instance.
# For details see https://developer.zendesk.com/api-reference/ticketing/tickets/custom_ticket_statuses/
# Following example creates a custom status with the category "open" and the agent label "In Progress".
# The example uses all possible attributes of the custom status resource.
# !! Warning: A custom status cannot be deleted once it has been created. It can only be deactivated.
# Deletion of the Resource will not delete the custom status in Zendesk.
resource "zendesk_custom_status" "example_status" {
  custom_status = {
    # Choose one of the categories: new, open, pending, hold, or solved
    status_category = "open"
    # Label that will be displayed to agents in the UI. Must be unique.
    agent_label = "In Progress"
    # Label that will be displayed to end users in the UI
    end_user_label = "We are on it!"
    # Description of the status for agents
    description = "This is an example progress status"
    # Description of the status for end users
    end_user_description = "Your request is being processed."
    # Whether the status is active or not. Not active would mean that the status is not available for selection in the UI
    active = true
  }

}