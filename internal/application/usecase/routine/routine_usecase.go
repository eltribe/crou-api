package routine

import (
	"crou-api/config"
	"crou-api/internal/application/inputport"
	"crou-api/internal/application/outputport"
	"crou-api/internal/domains"
	"crou-api/messages"
	"crou-api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RoutineUseCase struct {
	jwtProvider           *utils.JwtProvider
	routineDataOutputPort outputport.RoutineDataOutputPort
}

func NewRoutineUseCase(cnf *config.Config, routineDataOutputPort outputport.RoutineDataOutputPort) inputport.RoutineInputPort {
	return &RoutineUseCase{routineDataOutputPort: routineDataOutputPort, jwtProvider: utils.NewJwtProvider(cnf)}
}

func (svc RoutineUseCase) GetRoutines(c *fiber.Ctx) ([]*messages.RoutineResponse, error) {
	claims := svc.jwtProvider.GetClaims(c)
	routines, err := svc.routineDataOutputPort.GetRoutinesByUserID(uuid.MustParse(claims.Sub))
	if err != nil {
		return nil, err
	}

	routineResponses := make([]*messages.RoutineResponse, 0)
	for _, routine := range routines {
		routineResponses = append(routineResponses, &messages.RoutineResponse{
			ID:               routine.ID,
			Category:         routine.Category,
			RoutineType:      routine.RoutineType,
			Title:            routine.Title,
			When:             routine.When,
			TimeOfDay:        routine.TimeOfDay,
			DaysOfWeek:       routine.DaysOfWeek,
			IsNotification:   routine.IsNotification,
			NotificationTime: routine.NotificationTime,
		})
	}
	return routineResponses, nil
}

func (svc RoutineUseCase) CreateRoutine(c *fiber.Ctx, req messages.CreateRoutineRequest) (*messages.RoutineResponse, error) {
	claims := svc.jwtProvider.GetClaims(c)
	newRoutine := &domains.RoutineTemplate{
		UserId:           uuid.MustParse(claims.Sub),
		Category:         req.Category,
		RoutineType:      req.RoutineType,
		Title:            req.Title,
		When:             req.When,
		TimeOfDay:        req.TimeOfDay,
		DaysOfWeek:       req.DaysOfWeek,
		IsNotification:   req.IsNotification,
		NotificationTime: req.NotificationTime,
	}

	routine, err := svc.routineDataOutputPort.CreateRoutine(newRoutine)
	if err != nil {
		return nil, err
	}

	return &messages.RoutineResponse{
		ID:               routine.ID,
		Category:         routine.Category,
		RoutineType:      routine.RoutineType,
		Title:            routine.Title,
		When:             routine.When,
		TimeOfDay:        routine.TimeOfDay,
		DaysOfWeek:       routine.DaysOfWeek,
		IsNotification:   routine.IsNotification,
		NotificationTime: routine.NotificationTime,
	}, nil
}

func (svc RoutineUseCase) UpdateRoutine(c *fiber.Ctx, req messages.UpdateRoutineRequest) (*messages.RoutineResponse, error) {
	claims := svc.jwtProvider.GetClaims(c)
	updatedRoutine := &domains.RoutineTemplate{
		UserId:           uuid.MustParse(claims.Sub),
		Category:         req.Category,
		RoutineType:      req.RoutineType,
		Title:            req.Title,
		When:             req.When,
		TimeOfDay:        req.TimeOfDay,
		DaysOfWeek:       req.DaysOfWeek,
		IsNotification:   req.IsNotification,
		NotificationTime: req.NotificationTime,
	}

	routine, err := svc.routineDataOutputPort.UpdateRoutine(req.RoutineId, updatedRoutine)
	if err != nil {
		return nil, err
	}

	return &messages.RoutineResponse{
		ID:               routine.ID,
		Category:         routine.Category,
		RoutineType:      routine.RoutineType,
		Title:            routine.Title,
		When:             routine.When,
		TimeOfDay:        routine.TimeOfDay,
		DaysOfWeek:       routine.DaysOfWeek,
		IsNotification:   routine.IsNotification,
		NotificationTime: routine.NotificationTime,
	}, nil
}

func (svc RoutineUseCase) DeleteRoutine(c *fiber.Ctx, routineId uuid.UUID) error {
	err := svc.routineDataOutputPort.DeleteRoutine(routineId)
	if err != nil {
		return err
	}
	return nil
}
