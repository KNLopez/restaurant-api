package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/KNLopez/restaurant-api/internal/models"
	"github.com/google/uuid"
)

type RestaurantRepository struct {
	db *sql.DB
}

func NewRestaurantRepository(db *sql.DB) *RestaurantRepository {
	return &RestaurantRepository{db: db}
}

func (r *RestaurantRepository) Create(ctx context.Context, restaurant *models.Restaurant) error {
	query := `
		INSERT INTO restaurants (id, name, description, manager_id, address, phone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	now := time.Now()
	restaurant.ID = uuid.New()
	restaurant.CreatedAt = now
	restaurant.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		restaurant.ID,
		restaurant.Name,
		restaurant.Description,
		restaurant.ManagerID,
		restaurant.Address,
		restaurant.Phone,
		restaurant.CreatedAt,
		restaurant.UpdatedAt,
	)

	return err
}

func (r *RestaurantRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Restaurant, error) {
	query := `
		SELECT id, name, description, manager_id, address, phone, created_at, updated_at
		FROM restaurants
		WHERE id = $1
	`

	restaurant := &models.Restaurant{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&restaurant.ID,
		&restaurant.Name,
		&restaurant.Description,
		&restaurant.ManagerID,
		&restaurant.Address,
		&restaurant.Phone,
		&restaurant.CreatedAt,
		&restaurant.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return restaurant, nil
}

func (r *RestaurantRepository) GetByManagerID(ctx context.Context, managerID uuid.UUID) ([]*models.Restaurant, error) {
	query := `
		SELECT id, name, description, manager_id, address, phone, created_at, updated_at
		FROM restaurants
		WHERE manager_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, managerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restaurants []*models.Restaurant
	for rows.Next() {
		restaurant := &models.Restaurant{}
		err := rows.Scan(
			&restaurant.ID,
			&restaurant.Name,
			&restaurant.Description,
			&restaurant.ManagerID,
			&restaurant.Address,
			&restaurant.Phone,
			&restaurant.CreatedAt,
			&restaurant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		restaurants = append(restaurants, restaurant)
	}

	return restaurants, nil
}

func (r *RestaurantRepository) Update(ctx context.Context, restaurant *models.Restaurant) error {
	query := `
		UPDATE restaurants
		SET name = $1,
			description = $2,
			manager_id = $3,
			address = $4,
			phone = $5,
			updated_at = $6
		WHERE id = $7
	`

	restaurant.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		restaurant.Name,
		restaurant.Description,
		restaurant.ManagerID,
		restaurant.Address,
		restaurant.Phone,
		restaurant.UpdatedAt,
		restaurant.ID,
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

func (r *RestaurantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM restaurants WHERE id = $1`

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
