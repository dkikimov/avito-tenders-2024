package middlewares

import (
	"avito-tenders/internal/api/employee"
)

type Manager struct {
	empRepo employee.Repository
}

func NewManager(empRepo employee.Repository) *Manager {
	return &Manager{empRepo: empRepo}
}
