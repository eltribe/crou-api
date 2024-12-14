package outputport

import (
	"crou-api/internal/domains"
	"github.com/google/uuid"
)

type RoutineDataOutputPort interface {
	GetRoutinesByUserID(userID uuid.UUID) ([]*domains.RoutineTemplate, error)
	CreateRoutine(newRoutine *domains.RoutineTemplate) (*domains.RoutineTemplate, error)
	UpdateRoutine(id uuid.UUID, updatedRoutine *domains.RoutineTemplate) (*domains.RoutineTemplate, error)
	DeleteRoutine(id uuid.UUID) error

	WriteRoutineRecode(newRoutineRecord *domains.RoutineRecord) (*domains.RoutineRecord, error)
}
