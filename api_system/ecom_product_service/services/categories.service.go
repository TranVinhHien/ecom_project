package services

import (
	"context"
	"database/sql"
	"fmt"
	"mime/multipart"

	util_assets "github.com/TranVinhHien/ecom_product_service/assets/util"
	db "github.com/TranVinhHien/ecom_product_service/db/sqlc"
	assets_services "github.com/TranVinhHien/ecom_product_service/services/assets"
	services "github.com/TranVinhHien/ecom_product_service/services/entity"
)

func (s *service) GetCategoris(ctx context.Context, cate_id string) (map[string]interface{}, *assets_services.ServiceError) {

	caterd, err := s.redis.GetCategoryTree(ctx, cate_id)
	if err != nil {
		fmt.Println("Error GetCategories:", err)
		return nil, assets_services.NewError(400, err)
	}
	if len(caterd) == 0 {
		cates, err := s.repository.ListCategories(ctx)
		list_cate := make([]services.Categorys, 0)
		for _, cat := range cates {
			list_cate = append(list_cate, services.Categorys{
				CategoryID: cat.CategoryID,
				Name:       cat.Name,
				Parent:     services.Narg[string]{Valid: cat.Parent.Valid, Data: cat.Parent.String},
				Key:        cat.Key,
				Path:       cat.Path.String,
			})
		}
		// cates[0].
		if err != nil {
			fmt.Println("Error ListCategories:", err)
			return nil, assets_services.NewError(400, err)
		}
		err = s.redis.AddCategories(ctx, list_cate)
		if err != nil {
			fmt.Println("Error AddCategories:", err)
			return nil, assets_services.NewError(400, err)
		}
		caterd, err = s.redis.GetCategoryTree(ctx, cate_id)
		if err != nil {
			fmt.Println("Error GetCategories:", err)
			return nil, assets_services.NewError(400, err)
		}
	}
	// fmt.Println("\nGet all categories:")
	// for _, cat := range caterd {
	// 	fmt.Printf("Category: %s (%s)\n", cat.Name, cat.CategoryID)
	// 	if cat.Childs.Valid {
	// 		for _, child := range cat.Childs.Data {
	// 			fmt.Printf("  - Child: %s (%s)\n", child.Name, child.CategoryID)
	// 			if child.Childs.Valid {
	// 				for _, child2 := range child.Childs.Data {
	// 					fmt.Printf("	  - Child 2: %s (%s)\n", child2.Name, child2.CategoryID)

	// 				}
	// 			}
	// 		}
	// 	}
	// }

	result, err := assets_services.HideFields(caterd, "categories", "parent")
	if err != nil {
		fmt.Println("Error HideFields:", err)
		return nil, assets_services.NewError(400, err)
	}
	return result, nil
}

// Add a new category: create in DB, refresh redis cache, return category tree
func (s *service) AddCategory(ctx context.Context, userName string, cat services.Categorys, file *multipart.FileHeader) *assets_services.ServiceError {
	// create in repository
	if cat.Name == "" {
		return assets_services.NewError(400, fmt.Errorf("category name is required"))
	}
	if cat.Image.Valid {
		/// CALL API LUU TRONG SERVICE ẢNH
	}
	// create in repository
	key, err := assets_services.ConvertToSlug(cat.Name)
	if err != nil {
		return assets_services.NewError(400, fmt.Errorf("invalid category name: %v", err))
	}
	err = s.repository.CreateCategory(ctx, db.CreateCategoryParams{
		CategoryID: util_assets.RandomString(36),
		Name:       cat.Name,
		Parent:     sql.NullString{String: cat.Parent.Data, Valid: cat.Parent.Valid},
		Key:        key,
		Image:      sql.NullString{String: cat.Image.Data, Valid: cat.Image.Valid},
		// Path:       sql.NullString{String: key, Valid: true},
	})
	if err != nil {
		fmt.Println("Error CreateCategory:", err)
		return assets_services.NewError(400, err)
	}

	// refresh categories cache: list from DB -> convert -> add to redis
	cates, err := s.repository.ListCategories(ctx)
	if err != nil {
		fmt.Println("Error ListCategories after create:", err)
		return assets_services.NewError(400, err)
	}
	list_cate := make([]services.Categorys, 0, len(cates))
	for _, c := range cates {
		list_cate = append(list_cate, services.Categorys{
			CategoryID: c.CategoryID,
			Name:       c.Name,
			Parent:     services.Narg[string]{Valid: c.Parent.Valid, Data: c.Parent.String},
			Key:        c.Key,
			Path:       c.Path.String,
		})
	}
	if err = s.redis.AddCategories(ctx, list_cate); err != nil {
		fmt.Println("Error AddCategories after create:", err)
		return assets_services.NewError(400, err)
	}

	// return updated tree
	// tree, err := s.redis.GetCategoryTree(ctx, userName)
	// if err != nil {
	// 	fmt.Println("Error GetCategoryTree after create:", err)
	// 	return assets_services.NewError(400, err)
	// }
	// result, err := assets_services.HideFields(tree, "categories", "parent")
	// if err != nil {
	// 	fmt.Println("Error HideFields after create:", err)
	// 	return assets_services.NewError(400, err)
	// }
	// optionally include created item in response (created not used here)
	// _ = created
	return nil
}

// Update an existing category: update in DB, refresh redis cache, return category tree
func (s *service) UpdateCategory(ctx context.Context, userName string, cat services.Categorys, file *multipart.FileHeader) *assets_services.ServiceError {
	// update in repository
	if cat.CategoryID == "" {
		return assets_services.NewError(400, fmt.Errorf("category_id is required"))
	}
	if cat.Image.Valid {
		/// CALL API LUU TRONG SERVICE ẢNH
	}
	if cat.Name != "" {

		key, err := assets_services.ConvertToSlug(cat.Name)
		if err != nil {
			return assets_services.NewError(400, fmt.Errorf("invalid category name: %v", err))
		}
		cat.Key = key
	}

	if err := s.repository.UpdateCategory(ctx, db.UpdateCategoryParams{
		CategoryID: cat.CategoryID,
		Name:       sql.NullString{String: cat.Name, Valid: cat.Name != ""},
		Parent:     sql.NullString{String: cat.Parent.Data, Valid: cat.Parent.Valid},
		Key:        sql.NullString{String: cat.Key, Valid: cat.Key != ""},
		Image:      sql.NullString{String: cat.Image.Data, Valid: cat.Image.Valid},
		Path:       sql.NullString{String: cat.Key, Valid: cat.Key != ""},
	}); err != nil {
		fmt.Println("Error UpdateCategory:", err)
		return assets_services.NewError(400, err)
	}

	// refresh cache
	cates, err := s.repository.ListCategories(ctx)
	if err != nil {
		fmt.Println("Error ListCategories after update:", err)
		return assets_services.NewError(400, err)
	}
	list_cate := make([]services.Categorys, 0, len(cates))
	for _, c := range cates {
		list_cate = append(list_cate, services.Categorys{
			CategoryID: c.CategoryID,
			Name:       c.Name,
			Parent:     services.Narg[string]{Valid: c.Parent.Valid, Data: c.Parent.String},
			Key:        c.Key,
			Path:       c.Path.String,
		})
	}
	if err = s.redis.AddCategories(ctx, list_cate); err != nil {
		fmt.Println("Error AddCategories after update:", err)
		return assets_services.NewError(400, err)
	}
	return nil
}

// Delete a category: remove in DB, refresh redis cache, return category tree
func (s *service) DeleteCategory(ctx context.Context, userName string, categoryID string) *assets_services.ServiceError {
	// delete in repository
	if err := s.repository.DeleteCategory(ctx, categoryID); err != nil {
		fmt.Println("Error DeleteCategory:", err)
		return assets_services.NewError(400, err)
	}

	// refresh cache
	cates, err := s.repository.ListCategories(ctx)
	if err != nil {
		fmt.Println("Error ListCategories after delete:", err)
		return assets_services.NewError(400, err)
	}
	list_cate := make([]services.Categorys, 0, len(cates))
	for _, c := range cates {
		list_cate = append(list_cate, services.Categorys{
			CategoryID: c.CategoryID,
			Name:       c.Name,
			Parent:     services.Narg[string]{Valid: c.Parent.Valid, Data: c.Parent.String},
			Key:        c.Key,
			Path:       c.Path.String,
		})
	}
	if err = s.redis.AddCategories(ctx, list_cate); err != nil {
		fmt.Println("Error AddCategories after delete:", err)
		return assets_services.NewError(400, err)
	}

	return nil
}
