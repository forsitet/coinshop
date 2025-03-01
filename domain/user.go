package domain

type User struct {
	Username  string          `json:"username"`
	Inventory []InventoryItem `json:"inventory,omitempty"`
	ID        uint            `json:"id,omitempty"`
	Balance   int             `json:"balance"`
}
