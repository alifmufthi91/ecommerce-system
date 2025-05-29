package payload

import "github.com/google/uuid"

type UpdateWarehouseReq struct {
	ID     uuid.UUID `json:"-"`
	Status string    `json:"status" binding:"required,oneof=active inactive"`
}
