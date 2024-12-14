package messages

import (
	"crou-api/internal/domains"
	"github.com/google/uuid"
)

type CreateRoutineRequest struct {
	Category         domains.RoutineCategory `json:"category" validate:"required,oneof=SIGNATURE FAITH DAILY"`
	RoutineType      domains.RoutineType     `json:"routineType" validate:"required,oneof=BIBLE_TALK CHECK WRITE"`
	Title            string                  `json:"title" validate:"required"`
	When             string                  `json:"when" validate:"required"`
	TimeOfDay        domains.TimeOfDay       `json:"timeOfDay" validate:"required,oneof=MORNING AFTERNOON EVENING"`
	DaysOfWeek       []string                `json:"daysOfWeek" validate:"required,dive,oneof=MON TUE WED THU FRI SAT SUN"`
	IsNotification   bool                    `json:"isNotification"`
	NotificationTime *int32                  `json:"notificationTime"`
}

type UpdateRoutineRequest struct {
	RoutineId        uuid.UUID               `json:"routineId" validate:"required"`
	Category         domains.RoutineCategory `json:"category" validate:"required,oneof=SIGNATURE FAITH DAILY"`
	RoutineType      domains.RoutineType     `json:"routineType" validate:"required,oneof=BIBLE_TALK CHECK WRITE"`
	Title            string                  `json:"title" validate:"required"`
	When             string                  `json:"when" validate:"required"`
	TimeOfDay        domains.TimeOfDay       `json:"timeOfDay" validate:"required,oneof=MORNING AFTERNOON EVENING"`
	DaysOfWeek       []string                `json:"daysOfWeek" validate:"required,dive,oneof=MON TUE WED THU FRI SAT SUN"`
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
	DaysOfWeek       []string                `json:"daysOfWeek"`
	IsNotification   bool                    `json:"isNotification"`
	NotificationTime *int32                  `json:"notificationTime"`
}
