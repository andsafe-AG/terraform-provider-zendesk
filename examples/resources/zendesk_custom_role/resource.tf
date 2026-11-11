# Custom Agent Role to be used in your Zendesk instance. Requires the Enterprise plan or above.
# For details see https://developer.zendesk.com/api-reference/ticketing/account-configuration/custom_roles/
# Following example creates a custom role for partners, who can only make private comments on assigned tickets.
# All configuration attributes are optional. Attributes not set are left to the Zendesk defaults.
resource "zendesk_custom_role" "partner" {
  custom_role = {
    # Name of the custom role
    name = "Partner"
    # Description of the custom role
    description = "Can only make private comments on assigned tickets"
    configuration = {
      # What kind of tickets the agent can access: all, assigned-only, within-groups, within-groups-and-public-groups, within-organization
      ticket_access = "assigned-only"
      # What type of comments the agent can make: public, none
      ticket_comment_access = "none"
      # What the agent can do with end-user profiles: edit, edit-within-org, full, readonly
      end_user_profile_access = "readonly"
      # What the agent can do with views: full, manage-group, manage-personal, playonly, readonly
      view_access = "playonly"
      # What the agent can do with macros: full, manage-group, manage-personal, readonly
      macro_access    = "readonly"
      ticket_editing  = true
      ticket_deletion = false
      # Permissions (scopes) for custom objects: read, update, delete, create
      custom_objects = {
        shipment = {
          scopes = ["read"]
        }
      }
    }
  }
}
