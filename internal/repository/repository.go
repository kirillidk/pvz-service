package repository

import "database/sql"

type Repository struct {
	UserRepository      *UserRepository
	PVZRepository       *PVZRepository
	ReceptionRepository *ReceptionRepository
	ProductRepository   *ProductRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		UserRepository:      NewUserRepository(db),
		PVZRepository:       NewPVZRepository(db),
		ReceptionRepository: NewReceptionRepository(db),
		ProductRepository:   NewProductRepository(db),
	}
}
