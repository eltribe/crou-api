package gormadapter

import (
	"crou-api/config/database"
	"crou-api/internal/application/outputport"
	"crou-api/internal/domains"
	"github.com/google/uuid"
)

type RoutineGorm struct {
	database database.Persistent
}

func NewRoutineGorm(database database.Persistent) outputport.RoutineDataOutputPort {
	return &RoutineGorm{database: database}
}

func (r RoutineGorm) ListRoutinesByUserId(userID uuid.UUID) ([]*domains.Routine, error) {
	sql := r.database.DB()
	var routines []*domains.Routine
	result := sql.Model(&domains.Routine{}).Where("user_id = ?", userID).Find(&routines)
	if result.Error != nil {
		return nil, result.Error
	}
	return routines, nil
}

func (r RoutineGorm) GetRoutineById(id uuid.UUID) (*domains.Routine, error) {
	sql := r.database.DB()
	var routine domains.Routine
	result := sql.Model(&domains.Routine{}).Where("id = ?", id).First(&routine)
	if result.Error != nil {
		return nil, result.Error
	}
	return &routine, nil
}

func (r RoutineGorm) CreateRoutine(newRoutine *domains.Routine) (*domains.Routine, error) {
	sql := r.database.DB()
	result := sql.Create(newRoutine)
	if result.Error != nil {
		return nil, result.Error
	}
	return newRoutine, nil
}

func (r RoutineGorm) UpdateRoutine(id uuid.UUID, updatedRoutine *domains.Routine) (*domains.Routine, error) {
	sql := r.database.DB()
	result := sql.Model(&domains.Routine{}).Where("id = ?", id).Updates(updatedRoutine)
	if result.Error != nil {
		return nil, result.Error
	}
	return updatedRoutine, nil
}

func (r RoutineGorm) DeleteRoutine(id uuid.UUID) error {
	sql := r.database.DB()
	result := sql.Delete(&domains.Routine{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// bulk insert
func (r RoutineGorm) CreateRoutineSetByRoutines(y, m, d int, routines []*domains.Routine) ([]*domains.RoutineSet, error) {
	sql := r.database.DB()
	routineSetList := make([]*domains.RoutineSet, 0)
	for _, routine := range routines {
		routineSetList = append(routineSetList, &domains.RoutineSet{
			UserId:          routine.UserId,
			RoutineTemplate: routine.RoutineTemplate,
			Year:            y,
			Month:           m,
			Day:             d,
			RoutineId:       routine.ID,
		})
	}
	result := sql.CreateInBatches(routineSetList, len(routineSetList))
	if result.Error != nil {
		return nil, result.Error
	}
	return routineSetList, nil
}

func (r RoutineGorm) GetRoutineSetByRoutineIdAndDate(userId uuid.UUID, routineId uuid.UUID, Year int, Month int, Day int) (*domains.RoutineSet, error) {
	sql := r.database.DB()
	var routineSet domains.RoutineSet
	result := sql.Model(&domains.RoutineSet{}).Where("user_id = ? ANd routine_id = ? AND year = ? AND month = ? AND day = ?", userId, routineId, Year, Month, Day).First(&routineSet)
	if result.Error != nil {
		return nil, result.Error
	}
	return &routineSet, nil
}

func (r RoutineGorm) GetRoutineRecordBySetId(setId uuid.UUID) (*domains.RoutineRecord, error) {
	sql := r.database.DB()
	var routineRecord domains.RoutineRecord
	result := sql.Model(&domains.RoutineRecord{}).Where("routine_set_id = ?", setId).First(&routineRecord)
	if result.Error != nil {
		return nil, result.Error
	}
	return &routineRecord, nil
}

func (r RoutineGorm) WriteRoutineRecord(newRoutineRecord *domains.RoutineRecord) (*domains.RoutineRecord, error) {
	sql := r.database.DB()
	result := sql.Create(newRoutineRecord)
	if result.Error != nil {
		return nil, result.Error
	}
	return newRoutineRecord, nil
}

func (r RoutineGorm) DeleteRoutineRecord(routine *domains.RoutineRecord) error {
	sql := r.database.DB()
	result := sql.Delete(routine)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
