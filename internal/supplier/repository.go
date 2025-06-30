package supplier

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll(page, limit int) ([]Supplier, error)
	FindByID(id uint) (*Supplier, error)
	Create(supplier *Supplier) error
	Update(supplier *Supplier) error
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

// Create implements Repository.
func (r *repository) Create(supplier *Supplier) error {
	return r.db.Create(supplier).Error
}

// Delete implements Repository.
func (*repository) Delete(id uint) error {
	panic("unimplemented")
}

// FindAll implements Repository.
func (r *repository) FindAll(page, limit int) ([]Supplier, error) {
	var suppliers []Supplier
	offset := (page - 1) * limit
	if err := r.db.Limit(limit).Offset(offset).Find(&suppliers).Error; err != nil {
		return nil, err
	}
	return suppliers, nil
}

// FindByID implements Repository.
func (r *repository) FindByID(id uint) (*Supplier, error) {
	var supplier *Supplier
	if err := r.db.First(&supplier, id).Error; err != nil {
		return nil, err
	}
	return supplier, nil
}

// Update implements Repository.
func (r *repository) Update(supplier *Supplier) error {
	return r.db.Model(&supplier).Updates(supplier).Error
}
