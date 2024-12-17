package messages

import (
	"crou-api/internal/domains"
	"github.com/google/uuid"
	"time"
)

type RoutineTemplateDTO struct {
	Category    domains.RoutineCategory `json:"category" validate:"required,oneof=SIGNATURE FAITH DAILY"`
	RoutineType domains.RoutineType     `json:"routineType" validate:"required,oneof=BIBLE_TALK CHECK WRITE"`
	Title       string                  `json:"title" validate:"required"`
	When        string                  `json:"when" validate:"required"`
	TimeOfDay   domains.TimeOfDay       `json:"timeOfDay" validate:"required,oneof=MORNING AFTERNOON EVENING"`
}

type CreateRoutineRequest struct {
	RoutineTemplateDTO
	DaysOfWeek       []int32 `json:"daysOfWeek" validate:"required,dive,oneof=0 1 2 3 4 5 6"`
	IsNotification   bool    `json:"isNotification"`
	NotificationTime *int32  `json:"notificationTime"`
}

type UpdateRoutineRequest struct {
	RoutineId        uuid.UUID               `json:"routineId" validate:"required"`
	Category         domains.RoutineCategory `json:"category" validate:"required,oneof=SIGNATURE FAITH DAILY"`
	RoutineType      domains.RoutineType     `json:"routineType" validate:"required,oneof=BIBLE_TALK CHECK WRITE"`
	Title            string                  `json:"title" validate:"required"`
	When             string                  `json:"when" validate:"required"`
	TimeOfDay        domains.TimeOfDay       `json:"timeOfDay" validate:"required,oneof=MORNING AFTERNOON EVENING"`
	DaysOfWeek       []int32                 `json:"daysOfWeek" validate:"required,dive,oneof=0 1 2 3 4 5 6"`
	IsNotification   bool                    `json:"isNotification"`
	NotificationTime *int32                  `json:"notificationTime"`
}

type RoutineResponse struct {
	ID               uuid.UUID               `json:"id"`
	Category         domains.RoutineCategory `json:"category"`
	RoutineType      domains.RoutineType     `json:"routineType"`
	Title            string                  `json:"title"`
	When             string                  `json:"when"`
	TimeOfDay        domains.TimeOfDay       `json:"timeOfDay"`
	DaysOfWeek       []time.Weekday          `json:"daysOfWeek"`
	IsNotification   bool                    `json:"isNotification"`
	NotificationTime *int32                  `json:"notificationTime"`
}

type WriteRoutineRecordRequest struct {
	RoutineId     uuid.UUID `json:"routineId" validate:"required"`
	RecordContent *string   `json:"recordContent"`
	Year          int       `json:"year"`
	Month         int       `json:"month"`
	Day           int       `json:"day"`
}

type RollbackRoutineRecordRequest struct {
	RoutineId uuid.UUID `json:"routineId" validate:"required"`
	Year      int       `json:"year"`
	Month     int       `json:"month"`
	Day       int       `json:"day"`
}

type RoutineRecordResponse struct {
	IsRecord      bool    `json:"isRecord"`
	RecordContent *string `json:"recordContent"`
}

func ConvertToRoutineDTO(routine *domains.Routine) *RoutineResponse {
	return &RoutineResponse{
		ID:               routine.ID,
		Category:         routine.Category,
		RoutineType:      routine.RoutineType,
		Title:            routine.Title,
		When:             routine.When,
		TimeOfDay:        routine.TimeOfDay,
		DaysOfWeek:       routine.GetDaysOfWeek(),
		IsNotification:   routine.IsNotification,
		NotificationTime: routine.NotificationTime,
	}
}

func ConvertToRoutineRecordDTO(routineRecord *domains.RoutineRecord) *RoutineRecordResponse {
	return &RoutineRecordResponse{
		IsRecord:      routineRecord.IsRecord,
		RecordContent: routineRecord.RecordContent,
	}
}
