package supplier

type Service interface {
	GetAllSuppliers(page, limit int) ([]Supplier, error)
	GetSupplierByID(id uint) (*Supplier, error)
	CreateSupplier(supplier *Supplier) error
	UpdateSupplier(supplier *Supplier) error
}
type SupplierService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &SupplierService{repo}
}

func (s *SupplierService) CreateSupplier(supplier *Supplier) error {
	return s.repo.Create(supplier)
}

func (s *SupplierService) UpdateSupplier(supplier *Supplier) error {
	return s.repo.Update(supplier)
}

func (s *SupplierService) GetAllSuppliers(page, limit int) ([]Supplier, error) {
	return s.repo.FindAll(page, limit)
}

func (s *SupplierService) GetSupplierByID(id uint) (*Supplier, error) {
	return s.repo.FindByID(id)
}
