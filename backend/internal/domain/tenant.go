package domain

import "time"

type Tenant struct {
	ID                string    `json:"id" db:"id"`
	Name              string    `json:"name" db:"name"`
	Subdomain         string    `json:"subdomain" db:"subdomain"`
	Phone             string    `json:"phone" db:"phone"`
	Address           string    `json:"address" db:"address"`
	LogoURL           string    `json:"logo_url" db:"logo_url"`
	Plan              string    `json:"plan" db:"plan"`
	MidtransClientKey string    `json:"-" db:"midtrans_client_key"`
	MidtransServerKey string    `json:"-" db:"midtrans_server_key"`
	IsActive          bool      `json:"is_active" db:"is_active"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}
