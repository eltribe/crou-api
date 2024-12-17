package routine

import (
	"crou-api/config"
	"crou-api/errorcode"
	"crou-api/internal/application/inputport"
	"crou-api/internal/application/outputport"
	"crou-api/internal/domains"
	"crou-api/messages"
	"crou-api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DaysOfWeek []int32

type RoutineUseCase struct {
	jwtProvider       *utils.JwtProvider
	routineOutputPort outputport.RoutineDataOutputPort
}

func NewRoutineUseCase(cnf *config.Config, routineDataOutputPort outputport.RoutineDataOutputPort) inputport.RoutineInputPort {
	return &RoutineUseCase{routineOutputPort: routineDataOutputPort, jwtProvider: utils.NewJwtProvider(cnf)}
}

func (svc RoutineUseCase) GetRoutines(c *fiber.Ctx) ([]*messages.RoutineResponse, error) {
	claims := svc.jwtProvider.GetClaims(c)
	routines, err := svc.routineOutputPort.ListRoutinesByUserId(uuid.MustParse(claims.Sub))
	if err != nil {
		return nil, errorcode.ErrRoutineNotFound
	}

	routineResponses := make([]*messages.RoutineResponse, 0)
	for _, routine := range routines {
		routineResponses = append(routineResponses, messages.ConvertToRoutineDTO(routine))
	}
	return routineResponses, nil
}

func (svc RoutineUseCase) CreateRoutine(c *fiber.Ctx, req messages.CreateRoutineRequest) (*messages.RoutineResponse, error) {
	claims := svc.jwtProvider.GetClaims(c)
	newRoutine := &domains.Routine{
		UserId: uuid.MustParse(claims.Sub),
		RoutineTemplate: domains.RoutineTemplate{
			Category:    req.Category,
			RoutineType: req.RoutineType,
			Title:       req.Title,
			When:        req.When,
			TimeOfDay:   req.TimeOfDay,
		},
		DaysOfWeek:       req.DaysOfWeek,
		IsNotification:   req.IsNotification,
		NotificationTime: req.NotificationTime,
	}

	routine, err := svc.routineOutputPort.CreateRoutine(newRoutine)
	if err != nil {
		return nil, err
	}

	return messages.ConvertToRoutineDTO(routine), nil
}

func (svc RoutineUseCase) UpdateRoutine(c *fiber.Ctx, req messages.UpdateRoutineRequest) (*messages.RoutineResponse, error) {
	claims := svc.jwtProvider.GetClaims(c)
	updatedRoutine := &domains.Routine{
		UserId: uuid.MustParse(claims.Sub),
		RoutineTemplate: domains.RoutineTemplate{
			Category:    req.Category,
			RoutineType: req.RoutineType,
			Title:       req.Title,
			When:        req.When,
			TimeOfDay:   req.TimeOfDay,
		},
		DaysOfWeek:       req.DaysOfWeek,
		IsNotification:   req.IsNotification,
		NotificationTime: req.NotificationTime,
	}

	routine, err := svc.routineOutputPort.UpdateRoutine(req.RoutineId, updatedRoutine)
	if err != nil {
		return nil, err
	}

	return messages.ConvertToRoutineDTO(routine), nil
}

func (svc RoutineUseCase) DeleteRoutine(c *fiber.Ctx, routineId uuid.UUID) error {
	err := svc.routineOutputPort.DeleteRoutine(routineId)
	if err != nil {
		return err
	}
	return nil
}

func (svc RoutineUseCase) WriteRoutineRecord(c *fiber.Ctx, req messages.WriteRoutineRecordRequest) (*messages.RoutineRecordResponse, error) {
	claims := svc.jwtProvider.GetClaims(c)
	userId := uuid.MustParse(claims.Sub)
	routine, err := svc.routineOutputPort.GetRoutineById(req.RoutineId)
	if err != nil {
		return nil, errorcode.ErrRoutineNotFound
	}

	newRoutineRecord := &domains.RoutineRecord{
		RecordContent: req.RecordContent,
		IsRecord:      true,
	}

	// 루틴 세트 조회
	routineSet, err := svc.routineOutputPort.GetRoutineSetByRoutineIdAndDate(userId, routine.ID, req.Year, req.Month, req.Day)
	if err != nil { // 루틴 세트가 없으면 루틴 세트를 생성
		routines, err := svc.routineOutputPort.ListRoutinesByUserId(uuid.MustParse(claims.Sub))
		if err != nil {
			return nil, errorcode.ErrRoutineNotFound
		}
		// 루틴 세트 생성
		routineSetList, err := svc.routineOutputPort.CreateRoutineSetByRoutines(req.Year, req.Month, req.Day, routines)
		if err != nil {
			return nil, err
		}

		for _, rs := range routineSetList {
			if rs.RoutineId == routine.ID {
				newRoutineRecord.RoutineSetId = rs.ID
			}
		}

	} else {
		newRoutineRecord.RoutineSetId = routineSet.ID
	}

	alreadyRecord, _ := svc.routineOutputPort.GetRoutineRecordBySetId(newRoutineRecord.RoutineSetId)
	if alreadyRecord != nil {
		return nil, errorcode.ErrRoutineRecordAlreadyExist
	}

	routineRecord, err := svc.routineOutputPort.WriteRoutineRecord(newRoutineRecord)
	if err != nil {
		return nil, err
	}

	return messages.ConvertToRoutineRecordDTO(routineRecord), nil
}

func (svc RoutineUseCase) DeleteRoutineRecord(c *fiber.Ctx, req messages.RollbackRoutineRecordRequest) (*messages.RoutineRecordResponse, error) {
	claims := svc.jwtProvider.GetClaims(c)
	_, err := svc.routineOutputPort.GetRoutineById(uuid.MustParse(claims.Sub))
	if err != nil {
		return nil, errorcode.ErrRoutineNotFound
	}

	//svc.routineOutputPort.DeleteRoutineRecord(req.RoutineId, req.Year, req.Month, req.Day)
	return nil, nil
}
