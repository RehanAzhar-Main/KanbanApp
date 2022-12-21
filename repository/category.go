package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error)
	StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error)
	StoreManyCategory(ctx context.Context, categories []entity.Category) error
	GetCategoryByID(ctx context.Context, id int) (entity.Category, error)
	UpdateCategory(ctx context.Context, category *entity.Category) error
	DeleteCategory(ctx context.Context, id int) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error) {
	var listCategories []entity.Category

	rows, err := r.db.WithContext(ctx).Table("categories").
		Where("user_id = ?", id).Rows()

	if err != nil {
		fmt.Println("error di get catagories by user id")
		return []entity.Category{}, err
	}

	defer rows.Close()
	for rows.Next() {
		r.db.ScanRows(rows, &listCategories)
	}

	// categories not found
	if len(listCategories) == 0 {
		fmt.Println("error di get catagories by user id")
		return []entity.Category{}, nil
	}

	return listCategories, nil
}

func (r *categoryRepository) StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error) {
	if err := r.db.WithContext(ctx).Table("categories").
		Create(&category).Error; err != nil {
		return 0, err
	}

	return category.ID, nil
}

func (r *categoryRepository) StoreManyCategory(ctx context.Context, categories []entity.Category) error {
	if err := r.db.WithContext(ctx).Table("categories").
		Create(&categories).Error; err != nil {
		return err
	}

	return nil
}

func (r *categoryRepository) GetCategoryByID(ctx context.Context, id int) (entity.Category, error) {
	var category entity.Category

	// get category by id
	if err := r.db.WithContext(ctx).Table("categories").
		Where("id = ?", id).
		Find(&category).Error; err != nil {
		return entity.Category{}, err
	}

	// category not found
	if category == (entity.Category{}) {
		return entity.Category{}, nil
	}

	return category, nil
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, category *entity.Category) error {
	if err := r.db.WithContext(ctx).Table("categories").Where("id = ?", category.ID).
		Updates(&category).Error; err != nil {
		return err
	}

	return nil
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id int) error {
	var category entity.Category
	if err := r.db.WithContext(ctx).Delete(&category, id).Error; err != nil {
		return err
	}

	return nil
}
