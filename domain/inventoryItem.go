package domain

type InventoryItem struct {
	ID       uint   `json:"id,omitempty"`
	ItemType string `json:"item_type"`
	UserID   uint   `json:"user_id,omitempty"`
	Quantity int    `json:"quantity"`
}
