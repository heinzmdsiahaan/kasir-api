package repositories

import (
	"database/sql"
	"errors"
	models "kasir-api/model"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	query := "SELECT id, name FROM category"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categorys := make([]models.Category, 0)
	for rows.Next() {
		var p models.Category
		err := rows.Scan(&p.ID, &p.Name)
		if err != nil {
			return nil, err
		}
		categorys = append(categorys, p)
	}

	return categorys, nil
}

func (r *CategoryRepository) Create(category *models.Category) error {
	query := "INSERT INTO category (name) VALUES ($1) RETURNING id"
	err := r.db.QueryRow(query, category.Name).Scan(&category.ID)
	return err
}

// GetById gets a category by its ID
func (repo *CategoryRepository) GetByID(id int) (*models.Category, error) {
	query := "SELECT id, name  FROM category WHERE id = $1"

	var p models.Category
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name)
	if err == sql.ErrNoRows {
		return nil, errors.New("category not found")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// Update updates an existing category
func (repo *CategoryRepository) Update(category *models.Category) error {
	query := "UPDATE category SET name = $1 WHERE id = $2"
	result, err := repo.db.Exec(query, category.Name, category.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("category not found")
	}

	return nil
}

// Delete removes a category by its ID
func (repo *CategoryRepository) Delete(id int) error {
	query := "DELETE FROM category WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("category not found")
	}

	return err
}
