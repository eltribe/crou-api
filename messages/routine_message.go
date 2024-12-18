package messages

import (
	"crou-api/internal/domains"
	"github.com/google/uuid"
	"time"
)

type RoutineTemplateDTO struct {
	Category    domains.RoutineCategory `json:"category" validate:"required"`
	RoutineType domains.RoutineType     `json:"routineType" validate:"required"`
	Title       string                  `json:"title" validate:"required"`
	When        string                  `json:"when" validate:"required"`
	TimeOfDay   domains.TimeOfDay       `json:"timeOfDay" validate:"required,oneof=MORNING AFTERNOON EVENING"`
}

type CreateRoutineRequest struct {
	RoutineTemplateDTO
	UserId           uuid.UUID `json:"userId" swaggerignore:"true"`
	DaysOfWeek       []int32   `json:"daysOfWeek" validate:"required,dive,oneof=0 1 2 3 4 5 6"`
	IsNotification   bool      `json:"isNotification"`
	NotificationTime *int32    `json:"notificationTime"`
}

type UpdateRoutineRequest struct {
	RoutineId uuid.UUID `json:"routineId" swaggerignore:"true"`
	UserId    uuid.UUID `json:"userId" swaggerignore:"true"`
	RoutineTemplateDTO
	DaysOfWeek       []int32 `json:"daysOfWeek" validate:"required,dive,oneof=0 1 2 3 4 5 6"`
	IsNotification   bool    `json:"isNotification"`
	NotificationTime *int32  `json:"notificationTime"`
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
	RoutineId     uuid.UUID `json:"routineId" swaggerignore:"true"`
	UserId        uuid.UUID `json:"userId" swaggerignore:"true"`
	RecordContent *string   `json:"recordContent"`
	Year          int       `json:"year"`
	Month         int       `json:"month"`
	Day           int       `json:"day"`
}

type DeleteRoutineRecordRequest struct {
	UserId          uuid.UUID `json:"userId" swaggerignore:"true"`
	RoutineRecordId uuid.UUID `json:"routineRecordId" validate:"required"`
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
