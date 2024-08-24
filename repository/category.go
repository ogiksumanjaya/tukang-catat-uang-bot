package repository

import (
	"context"
	"database/sql"
	"github.com/ogiksumanjaya/entity"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (c *CategoryRepository) GetCategoryList(ctx context.Context, username string) ([]entity.Category, error) {
	var category []entity.Category
	query := `SELECT * FROM category WHERE username = $1`

	stmt, err := c.db.Prepare(query)
	if err != nil {
		return []entity.Category{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.Category
		err := rows.Scan(&item.ID, &item.Username, &item.Name)
		if err != nil {
			return nil, err
		}
		category = append(category, item)
	}

	return category, nil

}

func (c *CategoryRepository) GetCategoryByName(ctx context.Context, category entity.Category) (entity.Category, error) {
	query := `SELECT * FROM category WHERE username = $1 AND name = $2`

	stmt, err := c.db.Prepare(query)
	if err != nil {
		return entity.Category{}, err
	}
	defer stmt.Close()
	
	err = stmt.QueryRowContext(ctx, category.Username, category.Name).Scan(&category.ID, &category.Username, &category.Name)
	if err != nil {
		return entity.Category{}, err
	}

	return category, nil

}
