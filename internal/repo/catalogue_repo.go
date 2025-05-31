package repo

import (
	"context"
	"fmt"

	"github.com/godfreyowidi/simple-ecomm-demo/models"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatalogRepo struct {
	DB *pgxpool.Pool
}

func NewCatalogRepo(db *pgxpool.Pool) *CatalogRepo {
	return &CatalogRepo{DB: db}
}

func (r *CatalogRepo) GetProductCatalog(ctx context.Context) ([]models.ProductCatalog, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT 
			top_cat.name AS top_category_name,
			sub_cat.name AS sub_category_name,
			p.id, p.name, p.description, p.price
		FROM categories AS top_cat
		LEFT JOIN categories AS sub_cat ON sub_cat.parent_id = top_cat.id
		LEFT JOIN products p ON p.category_id = sub_cat.id
		WHERE top_cat.parent_id IS NULL
		ORDER BY top_cat.id, sub_cat.id, p.name
	`)
	if err != nil {
		return nil, fmt.Errorf("query catalog: %w", err)
	}
	defer rows.Close()

	catalogMap := map[string]*models.ProductCatalog{}

	for rows.Next() {
		var productID pgtype.Int4
		var productName *string
		var topCatName *string
		var subCatName *string
		var description *string
		var price pgtype.Float8

		err := rows.Scan(&topCatName, &subCatName, &productID, &productName, &description, &price)
		if err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		// Normalize category names
		topName := ""
		if topCatName != nil {
			topName = *topCatName
		}
		subName := ""
		if subCatName != nil {
			subName = *subCatName
		}

		// Get or create top category
		if _, ok := catalogMap[topName]; !ok {
			catalogMap[topName] = &models.ProductCatalog{
				TopCategoryName: topName,
			}
		}
		topCat := catalogMap[topName]

		// Get or create subcategory
		var subCat *models.ProductSubCategory
		for i := range topCat.SubCategories {
			if topCat.SubCategories[i].Name == subName {
				subCat = &topCat.SubCategories[i]
				break
			}
		}
		if subCat == nil {
			topCat.SubCategories = append(topCat.SubCategories, models.ProductSubCategory{Name: subName})
			subCat = &topCat.SubCategories[len(topCat.SubCategories)-1]
		}

		// Add product if present
		if productID.Valid && productName != nil {
			subCat.Products = append(subCat.Products, models.Product{
				ID:          int(productID.Int32),
				Name:        *productName,
				Description: description,
				Price:       price.Float64,
			})
		}
	}

	// Convert map to slice
	var result []models.ProductCatalog
	for _, c := range catalogMap {
		result = append(result, *c)
	}

	return result, nil
}
