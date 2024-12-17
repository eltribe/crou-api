package outputport

import (
	"crou-api/internal/domains"
	"github.com/google/uuid"
)

type RoutineDataOutputPort interface {
	ListRoutinesByUserId(userID uuid.UUID) ([]*domains.Routine, error)
	GetRoutineById(id uuid.UUID) (*domains.Routine, error)
	CreateRoutine(newRoutine *domains.Routine) (*domains.Routine, error)
	UpdateRoutine(id uuid.UUID, updatedRoutine *domains.Routine) (*domains.Routine, error)
	DeleteRoutine(id uuid.UUID) error

	CreateRoutineSetByRoutines(y, m, d int, routines []*domains.Routine) ([]*domains.RoutineSet, error)

	GetRoutineSetByRoutineIdAndDate(userId uuid.UUID, routineId uuid.UUID, Year int, Month int, Day int) (*domains.RoutineSet, error)
	GetRoutineRecordBySetId(setId uuid.UUID) (*domains.RoutineRecord, error)
	WriteRoutineRecord(newRoutineRecord *domains.RoutineRecord) (*domains.RoutineRecord, error)
	DeleteRoutineRecord(routine *domains.RoutineRecord) error
}
