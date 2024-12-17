package inputport

import (
	"crou-api/messages"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RoutineInputPort interface {
	GetRoutines(c *fiber.Ctx) ([]*messages.RoutineResponse, error)
	CreateRoutine(c *fiber.Ctx, req messages.CreateRoutineRequest) (*messages.RoutineResponse, error)
	UpdateRoutine(c *fiber.Ctx, req messages.UpdateRoutineRequest) (*messages.RoutineResponse, error)
	DeleteRoutine(c *fiber.Ctx, routineId uuid.UUID) error

	WriteRoutineRecord(c *fiber.Ctx, req messages.WriteRoutineRecordRequest) (*messages.RoutineRecordResponse, error)
	DeleteRoutineRecord(c *fiber.Ctx, req messages.RollbackRoutineRecordRequest) (*messages.RoutineRecordResponse, error)
}
