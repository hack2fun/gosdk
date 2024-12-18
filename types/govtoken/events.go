package types

// Govtoken module events
const (
	EventTypeMint         = "govtoken_mint"
	EventTypeOwnerChanged = "govtoken_owner_changed"

	AttributeKeyAmount    = "amount"
	AttributeKeyRecipient = "recipient"
	AttributeKeyOwner     = "owner"
	AttributeKeyOldOwner  = "old_owner"
	AttributeKeyNewOwner  = "new_owner"
)
