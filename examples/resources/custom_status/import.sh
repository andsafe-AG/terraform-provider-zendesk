# Custom Ticket Status can be imported by specifying the numeric identifier.
terraform import zendesk_custom_status.example 123
# It can also be imported by specifying the agent label.
terraform import zendesk_custom_status.example my-agent-label-for-custom-status