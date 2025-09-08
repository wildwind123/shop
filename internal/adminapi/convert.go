package adminapi

import (
	"github/wildwind123/shop/pkg/db"
	"github/wildwind123/shop/pkg/ogenapi"
)

// convertDBProductsToOgenAPI converts a slice of db.Product to ogenapi.Product
func convertDBProductsToOgenAPI(dbProducts []db.Product) []ogenapi.Product {
	ogenProducts := make([]ogenapi.Product, len(dbProducts))
	for i, dbProduct := range dbProducts {
		ogenProducts[i] = ogenapi.Product{
			ID:          dbProduct.ID,
			Name:        dbProduct.Name,
			Description: dbProduct.Description,
			CreatedAt:   dbProduct.CreatedAt,
			UpdatedAt:   dbProduct.UpdatedAt,
		}
	}
	return ogenProducts
}
