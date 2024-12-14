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

// 루틴 템플릿 (N개) -기록시-> 루틴 기록 (N개)
type RoutineTemplate struct {
	common.UUIDModel
	UserId      uuid.UUID `gorm:"index"`
	Category    RoutineCategory
	RoutineType RoutineType
	Title       string
	When        string    // 언제 할지?
	TimeOfDay   TimeOfDay // 오전,오후,저녁

	// 반복 주기 및 알림 설정
	DaysOfWeek       pq.StringArray `gorm:"type:text[]"` // 월,화,수,목,금,토,일
	IsNotification   bool
	NotificationTime *int32

	RoutineRecord *RoutineRecord
}

// 루틴 기록
type RoutineRecord struct {
	common.UUIDModel
	UserId      uuid.UUID `gorm:"index"`
	Category    RoutineCategory
	RoutineType RoutineType
	Title       string
	When        string    // 언제 할지?
	TimeOfDay   TimeOfDay // 오전,오후,저녁

	IsRecord        bool    // 기록했는지?
	RecordContent   *string `gorm:"type:text"`
	Year            int
	Month           int
	Day             int
	RecordDayOfWeek *time.Weekday

	RoutineTemplateId string
	RoutineTemplate   RoutineTemplate `gorm:"foreignKey:RoutineTemplateId;references:ID"`
}
