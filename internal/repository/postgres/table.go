package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/KNLopez/restaurant-api/internal/models"
	"github.com/google/uuid"
)

type TableRepository struct {
	db *sql.DB
}

func NewTableRepository(db *sql.DB) *TableRepository {
	return &TableRepository{db: db}
}

func (r *TableRepository) Create(ctx context.Context, table *models.Table) error {
	query := `
		INSERT INTO tables (
			id, restaurant_id, number, capacity,
			status, qr_code, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	now := time.Now()
	table.ID = uuid.New()
	table.CreatedAt = now
	table.UpdatedAt = now
	table.Status = models.TableStatusAvailable

	_, err := r.db.ExecContext(ctx, query,
		table.ID,
		table.RestaurantID,
		table.Number,
		table.Capacity,
		table.Status,
		table.QRCode,
		table.CreatedAt,
		table.UpdatedAt,
	)

	return err
}

func (r *TableRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Table, error) {
	query := `
		SELECT id, restaurant_id, number, capacity,
			   status, qr_code, created_at, updated_at
		FROM tables
		WHERE id = $1
	`

	table := &models.Table{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&table.ID,
		&table.RestaurantID,
		&table.Number,
		&table.Capacity,
		&table.Status,
		&table.QRCode,
		&table.CreatedAt,
		&table.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return table, nil
}

func (r *TableRepository) GetByRestaurantID(ctx context.Context, restaurantID uuid.UUID) ([]*models.Table, error) {
	query := `
		SELECT id, restaurant_id, number, capacity,
			   status, qr_code, created_at, updated_at
		FROM tables
		WHERE restaurant_id = $1
		ORDER BY number
	`

	rows, err := r.db.QueryContext(ctx, query, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []*models.Table
	for rows.Next() {
		table := &models.Table{}
		err := rows.Scan(
			&table.ID,
			&table.RestaurantID,
			&table.Number,
			&table.Capacity,
			&table.Status,
			&table.QRCode,
			&table.CreatedAt,
			&table.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

func (r *TableRepository) GetByQRCode(ctx context.Context, qrCode string) (*models.Table, error) {
	query := `
		SELECT id, restaurant_id, number, capacity,
			   status, qr_code, created_at, updated_at
		FROM tables
		WHERE qr_code = $1
	`

	table := &models.Table{}
	err := r.db.QueryRowContext(ctx, query, qrCode).Scan(
		&table.ID,
		&table.RestaurantID,
		&table.Number,
		&table.Capacity,
		&table.Status,
		&table.QRCode,
		&table.CreatedAt,
		&table.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return table, nil
}

func (r *TableRepository) Update(ctx context.Context, table *models.Table) error {
	query := `
		UPDATE tables
		SET number = $1,
			capacity = $2,
			status = $3,
			qr_code = $4,
			updated_at = $5
		WHERE id = $6 AND restaurant_id = $7
	`

	table.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		table.Number,
		table.Capacity,
		table.Status,
		table.QRCode,
		table.UpdatedAt,
		table.ID,
		table.RestaurantID,
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

func (r *TableRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.TableStatus) error {
	query := `
		UPDATE tables
		SET status = $1,
			updated_at = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query,
		status,
		time.Now(),
		id,
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

func (r *TableRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tables WHERE id = $1`

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
