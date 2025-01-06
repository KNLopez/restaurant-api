package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/KNLopez/restaurant-api/internal/models"
	"github.com/google/uuid"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create order
	query := `
		INSERT INTO orders (
			id, user_id, restaurant_id, status,
			total_amount, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	now := time.Now()
	order.ID = uuid.New()
	order.CreatedAt = now
	order.UpdatedAt = now
	order.Status = models.OrderStatusPending

	_, err = tx.ExecContext(ctx, query,
		order.ID,
		order.UserID,
		order.RestaurantID,
		order.Status,
		order.TotalAmount,
		order.CreatedAt,
		order.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Create order items
	itemQuery := `
		INSERT INTO order_items (
			id, order_id, menu_item_id,
			quantity, price
		) VALUES ($1, $2, $3, $4, $5)
	`

	for _, item := range order.Items {
		item.ID = uuid.New()
		item.OrderID = order.ID

		_, err = tx.ExecContext(ctx, itemQuery,
			item.ID,
			item.OrderID,
			item.MenuItemID,
			item.Quantity,
			item.Price,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *OrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	// Get order
	orderQuery := `
		SELECT id, user_id, restaurant_id, status,
			   total_amount, created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	order := &models.Order{}
	err := r.db.QueryRowContext(ctx, orderQuery, id).Scan(
		&order.ID,
		&order.UserID,
		&order.RestaurantID,
		&order.Status,
		&order.TotalAmount,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Get order items
	itemsQuery := `
		SELECT id, order_id, menu_item_id, quantity, price
		FROM order_items
		WHERE order_id = $1
	`

	rows, err := r.db.QueryContext(ctx, itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &models.OrderItem{}
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.MenuItemID,
			&item.Quantity,
			&item.Price,
		)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, *item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Order, error) {
	query := `
		SELECT id, user_id, restaurant_id, status,
			   total_amount, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		order := &models.Order{}
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.RestaurantID,
			&order.Status,
			&order.TotalAmount,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Get items for each order
	for _, order := range orders {
		itemsQuery := `
			SELECT id, order_id, menu_item_id, quantity, price
			FROM order_items
			WHERE order_id = $1
		`

		rows, err := r.db.QueryContext(ctx, itemsQuery, order.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			item := &models.OrderItem{}
			err := rows.Scan(
				&item.ID,
				&item.OrderID,
				&item.MenuItemID,
				&item.Quantity,
				&item.Price,
			)
			if err != nil {
				return nil, err
			}
			order.Items = append(order.Items, *item)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	return orders, nil
}

func (r *OrderRepository) GetByRestaurantID(ctx context.Context, restaurantID uuid.UUID) ([]*models.Order, error) {
	query := `
		SELECT id, user_id, restaurant_id, status,
			   total_amount, created_at, updated_at
		FROM orders
		WHERE restaurant_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		order := &models.Order{}
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.RestaurantID,
			&order.Status,
			&order.TotalAmount,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) Update(ctx context.Context, order *models.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update order
	query := `
		UPDATE orders
		SET status = $1,
			total_amount = $2,
			updated_at = $3
		WHERE id = $4
	`

	order.UpdatedAt = time.Now()

	result, err := tx.ExecContext(ctx, query,
		order.Status,
		order.TotalAmount,
		order.UpdatedAt,
		order.ID,
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

	// Delete existing items
	_, err = tx.ExecContext(ctx, "DELETE FROM order_items WHERE order_id = $1", order.ID)
	if err != nil {
		return err
	}

	// Insert updated items
	itemQuery := `
		INSERT INTO order_items (
			id, order_id, menu_item_id,
			quantity, price
		) VALUES ($1, $2, $3, $4, $5)
	`

	for _, item := range order.Items {
		if item.ID == uuid.Nil {
			item.ID = uuid.New()
		}
		item.OrderID = order.ID

		_, err = tx.ExecContext(ctx, itemQuery,
			item.ID,
			item.OrderID,
			item.MenuItemID,
			item.Quantity,
			item.Price,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.OrderStatus) error {
	query := `
		UPDATE orders
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

func (r *OrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete order items first
	_, err = tx.ExecContext(ctx, "DELETE FROM order_items WHERE order_id = $1", id)
	if err != nil {
		return err
	}

	// Delete order
	result, err := tx.ExecContext(ctx, "DELETE FROM orders WHERE id = $1", id)
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

	return tx.Commit()
}
