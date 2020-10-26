//go:generate reform
package gotrix

import "time"

//reform:catalogs
type Catalog struct {
	CatalogID   string    `json:"catalog_id" reform:"catalog_id,pk"`
	CatalogName string    `json:"catalog_name" reform:"catalog_name"`
	CreatedAt   time.Time `json:"created_at" reform:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" reform:"updated_at"`
}

//reform:catalog_items
type CatalogItem struct {
	ItemID    string    `json:"item_id" reform:"item_id,pk"`
	CatalogID string    `json:"catalog_id" reform:"catalog_id"`
	ItemName  string    `json:"item_name" reform:"item_name"`
	CreatedAt time.Time `json:"created_at" reform:"created_at"`
	UpdatedAt time.Time `json:"updated_at" reform:"updated_at"`
}

//reform:catalog_item_properties
type CatalogItemProperty struct {
	PropertyID    string    `json:"property_id" reform:"property_id,pk"`
	ItemID        string    `json:"item_id" reform:"item_id"`
	CatalogID     string    `json:"catalog_id" reform:"catalog_id"`
	PropertyKey   string    `json:"property_key" reform:"property_key"`
	PropertyValue string    `json:"property_value" reform:"property_value"`
	CreatedAt     time.Time `json:"created_at" reform:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" reform:"updated_at"`
}
