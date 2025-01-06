package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/KNLopez/restaurant-api/internal/models"
	"github.com/google/uuid"
)

type MenuRepository struct {
	db *sql.DB
}

func NewMenuRepository(db *sql.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

func (r *MenuRepository) Create(ctx context.Context, item *models.MenuItem) error {
	query := `
		INSERT INTO menu_items (
			id, restaurant_id, name, description, 
			price, category_id, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	now := time.Now()
	item.ID = uuid.New()
	item.CreatedAt = now
	item.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		item.ID,
		item.RestaurantID,
		item.Name,
		item.Description,
		item.Price,
		item.CategoryID,
		item.CreatedAt,
		item.UpdatedAt,
	)

	return err
}

func (r *MenuRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.MenuItem, error) {
	query := `
		SELECT id, restaurant_id, name, description, 
			   price, category_id, created_at, updated_at
		FROM menu_items
		WHERE id = $1
	`

	item := &models.MenuItem{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.RestaurantID,
		&item.Name,
		&item.Description,
		&item.Price,
		&item.CategoryID,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *MenuRepository) List(ctx context.Context, restaurantID uuid.UUID) ([]*models.MenuItem, error) {
	query := `
		SELECT id, restaurant_id, name, description, 
			   price, category_id, created_at, updated_at
		FROM menu_items
		WHERE restaurant_id = $1
		ORDER BY category_id, name
	`

	rows, err := r.db.QueryContext(ctx, query, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.MenuItem
	for rows.Next() {
		item := &models.MenuItem{}
		err := rows.Scan(
			&item.ID,
			&item.RestaurantID,
			&item.Name,
			&item.Description,
			&item.Price,
			&item.CategoryID,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *MenuRepository) Update(ctx context.Context, item *models.MenuItem) error {
	query := `
		UPDATE menu_items
		SET name = $1,
			description = $2,
			price = $3,
			category_id = $4,
			updated_at = $5
		WHERE id = $6 AND restaurant_id = $7
	`

	item.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		item.Name,
		item.Description,
		item.Price,
		item.CategoryID,
		item.UpdatedAt,
		item.ID,
		item.RestaurantID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *MenuRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM menu_items WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
