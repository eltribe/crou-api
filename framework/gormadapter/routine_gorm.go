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

func (r RoutineGorm) GetRoutinesByUserID(userID uuid.UUID) ([]*domains.RoutineTemplate, error) {
	sql := r.database.DB()
	var routines []*domains.RoutineTemplate
	result := sql.Model(&domains.RoutineTemplate{}).Preload("RoutineRecord").Where("user_id = ?", userID).Find(&routines)
	if result.Error != nil {
		return nil, result.Error
	}
	return routines, nil
}

func (r RoutineGorm) CreateRoutine(newRoutine *domains.RoutineTemplate) (*domains.RoutineTemplate, error) {
	sql := r.database.DB()
	result := sql.Create(newRoutine)
	if result.Error != nil {
		return nil, result.Error
	}
	return newRoutine, nil
}

func (r RoutineGorm) UpdateRoutine(id uuid.UUID, updatedRoutine *domains.RoutineTemplate) (*domains.RoutineTemplate, error) {
	sql := r.database.DB()
	result := sql.Model(&domains.RoutineTemplate{}).Where("id = ?", id).Updates(updatedRoutine)
	if result.Error != nil {
		return nil, result.Error
	}
	return updatedRoutine, nil
}

func (r RoutineGorm) DeleteRoutine(id uuid.UUID) error {
	sql := r.database.DB()
	result := sql.Delete(&domains.RoutineTemplate{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r RoutineGorm) WriteRoutineRecode(newRoutineRecord *domains.RoutineRecord) (*domains.RoutineRecord, error) {
	sql := r.database.DB()
	result := sql.Create(newRoutineRecord)
	if result.Error != nil {
		return nil, result.Error
	}
	return newRoutineRecord, nil
}
