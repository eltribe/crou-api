package domains

import (
	"crou-api/common"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

type RoutineCategory string

const (
	SIGNATURE RoutineCategory = "SIGNATURE"
	FAITH     RoutineCategory = "FAITH"
	DAILY     RoutineCategory = "DAILY"
)

type RoutineType string

// 기록형, 체크형, 바이블톡
const (
	BIBLE_TALK RoutineType = "BIBLE_TALK"
	CHECK      RoutineType = "CHECK"
	WRITE      RoutineType = "WRITE"
)

type TimeOfDay string

const (
	MORNING   TimeOfDay = "MORNING"
	AFTERNOON TimeOfDay = "AFTERNOON"
	EVENING   TimeOfDay = "EVENING"
)

type RoutineTemplate struct {
	Category    RoutineCategory `json:"category" validate:"required,oneof=SIGNATURE FAITH DAILY"`
	RoutineType RoutineType     `json:"routineType" validate:"required,oneof=BIBLE_TALK CHECK WRITE"`
	Title       string          `json:"title" validate:"required"`
	When        string          `json:"when" validate:"required"`
	TimeOfDay   TimeOfDay       `json:"timeOfDay" validate:"required,oneof=MORNING AFTERNOON EVENING"`
}

// 루틴 템플릿(N) -> 루틴세트(M) -> 루틴 기록(1)
type Routine struct {
	common.UUIDModel
	UserId uuid.UUID `gorm:"index"`
	RoutineTemplate

	// 반복 주기 및 알림 설정
	DaysOfWeek       pq.Int32Array `gorm:"type:integer[]"`
	IsNotification   bool
	NotificationTime *int32
}

func (rt Routine) GetDaysOfWeek() []time.Weekday {
	daysOfWeek := make([]time.Weekday, 0)
	for _, day := range rt.DaysOfWeek {
		daysOfWeek = append(daysOfWeek, time.Weekday(day))
	}
	return daysOfWeek
}

type RoutineSet struct {
	common.UUIDModel
	UserId uuid.UUID
	RoutineTemplate

	Year  int `gorm:"index:record_date"`
	Month int `gorm:"index:record_date"`
	Day   int `gorm:"index:record_date"`

	RoutineId uuid.UUID
	Routine   Routine `gorm:"foreignKey:RoutineId;references:ID"`
}

// 루틴 기록
type RoutineRecord struct {
	RoutineSetId uuid.UUID `gorm:"type:uuid;primaryKey;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	IsRecord      bool    // 기록했는지?
	RecordContent *string `gorm:"type:text"`

	RoutineSet RoutineSet `gorm:"foreignKey:RoutineSetId;references:ID"`
}
