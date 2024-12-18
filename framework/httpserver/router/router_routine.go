package router

import (
	"crou-api/config"
	app "crou-api/internal"
	"crou-api/messages"
	"crou-api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func routine(
	conf *config.Config,
	route *fiber.App,
	stx *app.InputPortProvider,
) {
	jwtProvider := utils.NewJwtProvider(conf)

	route.Post("/routine", func(c *fiber.Ctx) error {
		req := messages.CreateRoutineRequest{}
		if err := bodyValidator(c, &req); err != nil {
			return err
		}

		req.UserId = getUserId(c, jwtProvider)
		result, err := stx.RoutineInputPort.CreateRoutine(c, req)
		return response(c, result, err)
	})

	route.Get("/routine", func(c *fiber.Ctx) error {
		result, err := stx.RoutineInputPort.GetRoutines(c, getUserId(c, jwtProvider))
		return response(c, result, err)
	})

	route.Put("/routine/:routineId", func(c *fiber.Ctx) error {
		req := messages.UpdateRoutineRequest{}
		if err := bodyValidator(c, &req); err != nil {
			return err
		}
		req.RoutineId = uuid.MustParse(c.Params("routineId"))
		req.UserId = getUserId(c, jwtProvider)
		result, err := stx.RoutineInputPort.UpdateRoutine(c, req)
		return response(c, result, err)
	})

	route.Delete("/routine/:routineId", func(c *fiber.Ctx) error {
		routineId := c.Params("routineId")
		err := stx.RoutineInputPort.DeleteRoutine(c, uuid.MustParse(routineId))
		return response(c, nil, err)
	})

	route.Post("/routine/:routineId/record", func(c *fiber.Ctx) error {
		req := messages.WriteRoutineRecordRequest{}
		if err := bodyValidator(c, &req); err != nil {
			return err
		}
		req.RoutineId = uuid.MustParse(c.Params("routineId"))
		result, err := stx.RoutineInputPort.WriteRoutineRecord(c, req)
		return response(c, result, err)
	})

	route.Delete("/routine/record/:routineRecordId", func(c *fiber.Ctx) error {
		routineRecordId := c.Params("routineRecordId")
		err := stx.RoutineInputPort.DeleteRoutineRecord(c, messages.DeleteRoutineRecordRequest{
			RoutineRecordId: uuid.MustParse(routineRecordId),
			UserId:          getUserId(c, jwtProvider),
		})
		return response(c, nil, err)
	})
}
