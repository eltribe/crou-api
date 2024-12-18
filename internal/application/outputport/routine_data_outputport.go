package outputport

import (
	"crou-api/internal/domains"
	"github.com/google/uuid"
)

type RoutineDataOutputPort interface {

	// Routine
	ListRoutinesByUserId(userID uuid.UUID) ([]*domains.Routine, error)
	GetRoutineById(id uuid.UUID) (*domains.Routine, error)
	CreateRoutine(newRoutine *domains.Routine) (*domains.Routine, error)
	UpdateRoutine(id uuid.UUID, updatedRoutine *domains.Routine) (*domains.Routine, error)
	DeleteRoutine(id uuid.UUID) error

	// RoutineSet
	CreateRoutineSetByRoutines(y, m, d int, routines []*domains.Routine) ([]*domains.RoutineSet, error)
	GetRoutineSetByRoutineIdAndDate(userId uuid.UUID, routineId uuid.UUID, y, m, d int) (*domains.RoutineSet, error)

	// RoutineRecord
	WriteRoutineRecord(record *domains.RoutineRecord) (*domains.RoutineRecord, error)
	GetRoutineRecordBySetId(setId uuid.UUID) (*domains.RoutineRecord, error)
	GetRoutineRecordById(recordId uuid.UUID) (*domains.RoutineRecord, error)
	DeleteRoutineRecord(recordId uuid.UUID) error
}
