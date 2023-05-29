package db

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"github.com/Tanrungthip/mecari-build-hackathon-2023/backend/domain"
	"github.com/labstack/echo/v4"
)

type UserRepository interface {
	AddUser(ctx context.Context, user domain.User) (int64, error)
	GetUser(ctx context.Context, id int64) (domain.User, error)
	UpdateBalance(ctx context.Context, id int64, balance int64) error
}

type UserDBRepository struct {
	*sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &UserDBRepository{DB: db}
}

func (r *UserDBRepository) AddUser(ctx context.Context, user domain.User) (int64, error) {
	tx, err := r.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		fmt.Sprintf("failed to begin DB: %s\n", err)
		return 0, err
	}

	if _, err := tx.ExecContext(ctx, "INSERT INTO users (name, password) VALUES (?, ?)", user.Name, user.Password); err != nil {
		tx.Rollback()
		return 0, echo.NewHTTPError(http.StatusConflict, err)
	} else {
		tx.Commit()
	}

	// TODO: if other insert query is executed at the same time, it might return wrong id
	// http.StatusConflict(409) 既に同じIDがあった場合
	row := r.QueryRowContext(ctx, "SELECT id FROM users WHERE rowid = LAST_INSERT_ROWID()")

	var id int64
	return id, row.Scan(&id)
}

func (r *UserDBRepository) GetUser(ctx context.Context, id int64) (domain.User, error) {
	row := r.QueryRowContext(ctx, "SELECT * FROM users WHERE id = ?", id)

	var user domain.User
	return user, row.Scan(&user.ID, &user.Name, &user.Password, &user.Balance)
}

func (r *UserDBRepository) UpdateBalance(ctx context.Context, id int64, balance int64) error {
	tx, err := r.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		fmt.Sprintf("failed to begin DB: %s\n", err)
		// log.Fatal(err)
		return err
	}

	if _, err := r.ExecContext(ctx, "UPDATE users SET balance = ? WHERE id = ?", balance, id); err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusConflict, err)
	} else {
		tx.Commit()
	}
	return nil
}

type ItemRepository interface {
	AddItem(ctx context.Context, item domain.Item) (domain.Item, error)
	GetItem(ctx context.Context, id int32) (domain.Item, error)
	GetItemImage(ctx context.Context, id int32) ([]byte, error)
	GetOnSaleItems(ctx context.Context) ([]domain.Item, error)
	GetItemsByUserID(ctx context.Context, userID int64) ([]domain.Item, error)
	GetCategory(ctx context.Context, id int64) (domain.Category, error)
	GetCategories(ctx context.Context) ([]domain.Category, error)
	SearchItem(ctx context.Context, name string) ([]domain.Item, error)
	UpdateItem(ctx context.Context, id int32, item domain.Item) error
	UpdateItemStatus(ctx context.Context, id int32, status domain.ItemStatus) error
}

type ItemDBRepository struct {
	*sql.DB
}

func NewItemRepository(db *sql.DB) ItemRepository {
	return &ItemDBRepository{DB: db}
}

func (r *ItemDBRepository) AddItem(ctx context.Context, item domain.Item) (domain.Item, error) {
	tx, err := r.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		fmt.Sprintf("failed to begin DB: %s\n", err)
		// log.Fatal(err)
		return domain.Item{}, err
	}

	if _, err := tx.ExecContext(ctx, "INSERT INTO items (name, price, description, category_id, seller_id, image, status) VALUES (?, ?, ?, ?, ?, ?, ?)", item.Name, item.Price, item.Description, item.CategoryID, item.UserID, item.Image, item.Status); err != nil {
		tx.Rollback()
		return domain.Item{}, echo.NewHTTPError(http.StatusConflict, err)
	} else {
		tx.Commit()
	}

	// TODO: if other insert query is executed at the same time, it might return wrong id
	// http.StatusConflict(409) 既に同じIDがあった場合
	row := r.QueryRowContext(ctx, "SELECT * FROM items WHERE rowid = LAST_INSERT_ROWID()")

	var res domain.Item
	return res, row.Scan(&res.ID, &res.Name, &res.Price, &res.Description, &res.CategoryID, &res.UserID, &res.Image, &res.Status, &res.CreatedAt, &res.UpdatedAt)
}

func (r *ItemDBRepository) UpdateItem(ctx context.Context, id int32, item domain.Item) error {
	tx, err := r.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		fmt.Sprintf("failed to begin DB: %s\n", err)
		// log.Fatal(err)
		return err
	}

	if _, err := r.ExecContext(ctx, "UPDATE items SET name=?, price=?, description=?, category_id=? WHERE id=?", item.Name, item.Price, item.Description, item.CategoryID, id); err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusConflict, err)
	} else {
		tx.Commit()
	}
	return nil
}

func (r *ItemDBRepository) GetItem(ctx context.Context, id int32) (domain.Item, error) {
	row := r.QueryRowContext(ctx, "SELECT * FROM items WHERE id = ?", id)

	var item domain.Item
	return item, row.Scan(&item.ID, &item.Name, &item.Price, &item.Description, &item.CategoryID, &item.UserID, &item.Image, &item.Status, &item.CreatedAt, &item.UpdatedAt)
}

func (r *ItemDBRepository) SearchItem(ctx context.Context, name string) ([]domain.Item, error) {
	name = "%" + name + "%"
	rows, err := r.QueryContext(ctx, "SELECT * FROM items WHERE name LIKE ?", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var itemlist []domain.Item
	for rows.Next() {
		var item domain.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Description, &item.CategoryID, &item.UserID, &item.Image, &item.Status, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		itemlist = append(itemlist, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return itemlist, nil
}

func (r *ItemDBRepository) GetItemImage(ctx context.Context, id int32) ([]byte, error) {
	row := r.QueryRowContext(ctx, "SELECT image FROM items WHERE id = ?", id)
	var image []byte
	return image, row.Scan(&image)
}

func (r *ItemDBRepository) GetOnSaleItems(ctx context.Context) ([]domain.Item, error) {
	rows, err := r.QueryContext(ctx, "SELECT * FROM items WHERE status = ? ORDER BY updated_at desc", domain.ItemStatusOnSale)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.Item
	for rows.Next() {
		var item domain.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Description, &item.CategoryID, &item.UserID, &item.Image, &item.Status, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ItemDBRepository) GetItemsByUserID(ctx context.Context, userID int64) ([]domain.Item, error) {
	rows, err := r.QueryContext(ctx, "SELECT * FROM items WHERE seller_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.Item
	for rows.Next() {
		var item domain.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Description, &item.CategoryID, &item.UserID, &item.Image, &item.Status, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ItemDBRepository) UpdateItemStatus(ctx context.Context, id int32, status domain.ItemStatus) error {
	if _, err := r.ExecContext(ctx, "UPDATE items SET status = ? WHERE id = ?", status, id); err != nil {
		return err
	}
	return nil
}

func (r *ItemDBRepository) GetCategory(ctx context.Context, id int64) (domain.Category, error) {
	row := r.QueryRowContext(ctx, "SELECT * FROM category WHERE id = ?", id)

	var cat domain.Category
	return cat, row.Scan(&cat.ID, &cat.Name)
}

func (r *ItemDBRepository) GetCategories(ctx context.Context) ([]domain.Category, error) {
	rows, err := r.QueryContext(ctx, "SELECT * FROM category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []domain.Category
	for rows.Next() {
		var cat domain.Category
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			return nil, err
		}
		cats = append(cats, cat)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cats, nil
}
