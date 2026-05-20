package entity

import "time"

type Queue struct {
	ID           string     `json:"id"`
	QueueNumber  int        `json:"queue_number"`
	CustomerID   string     `json:"customer_id"`
	CustomerName string     `json:"customer_name"`
	BarberID     *string    `json:"barber_id,omitempty"`
	ServiceID    string     `json:"service_id"`
	ServiceName  string     `json:"service_name"`
	Status       string     `json:"status"`
	Position     int        `json:"position"`
	CreatedAt    time.Time  `json:"created_at"`
	CalledAt     *time.Time `json:"called_at,omitempty"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
}

type QueueStatus string

const (
	QueueStatusWaiting    QueueStatus = "waiting"
	QueueStatusCalled     QueueStatus = "called"
	QueueStatusInProgress QueueStatus = "in_progress"
	QueueStatusCompleted  QueueStatus = "completed"
	QueueStatusCanceled   QueueStatus = "canceled"
	QueueStatusSkipped    QueueStatus = "skipped"
)

type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleAdmin    UserRole = "admin"
	RoleBarber   UserRole = "barber"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Barber struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Specialty string    `json:"specialty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Service struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Duration    int     `json:"duration"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
