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

type RoutineUseCase struct {
	jwtProvider       *utils.JwtProvider
	routineOutputPort outputport.RoutineDataOutputPort
}

func NewRoutineUseCase(cnf *config.Config, routineDataOutputPort outputport.RoutineDataOutputPort) inputport.RoutineInputPort {
	return &RoutineUseCase{routineOutputPort: routineDataOutputPort, jwtProvider: utils.NewJwtProvider(cnf)}
}

// GetRoutines godoc
//
// @Summary		Get Routines API
// @Description	Get routines by user ID.
// @Accept		json
// @Produce		json
// @Tags 		데일리 루틴
// @Success		200	{object}	[]messages.RoutineResponse
// @Failure		409	{object}	errorcode.UseCaseError
// @Router		/v1/routines [get]
func (svc RoutineUseCase) GetRoutines(c *fiber.Ctx, userId uuid.UUID) ([]*messages.RoutineResponse, error) {
	routines, err := svc.routineOutputPort.ListRoutinesByUserId(userId)
	if err != nil {
		return nil, errorcode.ErrRoutineNotFound
	}

	routineResponses := make([]*messages.RoutineResponse, 0)
	for _, routine := range routines {
		routineResponses = append(routineResponses, messages.ConvertToRoutineDTO(routine))
	}
	return routineResponses, nil
}

// CreateRoutine godoc
//
// @Summary		Create Routine API
// @Description	Create a new routine.
// @Accept		json
// @Produce		json
// @Tags 		데일리 루틴
// @Param		messages.CreateRoutineRequest	body	messages.CreateRoutineRequest	true	"Create Routine Request"
// @Success		200	{object}	messages.RoutineResponse
// @Failure		409	{object}	errorcode.UseCaseError
// @Router		/v1/routines [post]
func (svc RoutineUseCase) CreateRoutine(c *fiber.Ctx, req messages.CreateRoutineRequest) (*messages.RoutineResponse, error) {
	newRoutine := &domains.Routine{
		UserId: req.UserId,
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

// UpdateRoutine godoc
//
// @Summary		Update Routine API
// @Description	Update an existing routine.
// @Accept		json
// @Produce		json
// @Tags 		데일리 루틴
// @Param		id	path	string	true	"Routine ID"
// @Param		messages.UpdateRoutineRequest	body	messages.UpdateRoutineRequest	true	"Update Routine Request"
// @Success		200	{object}	messages.RoutineResponse
// @Failure		409	{object}	errorcode.UseCaseError
// @Router		/v1/routines/{id} [put]
func (svc RoutineUseCase) UpdateRoutine(c *fiber.Ctx, req messages.UpdateRoutineRequest) (*messages.RoutineResponse, error) {
	updatedRoutine := &domains.Routine{
		UserId: req.UserId,
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

// DeleteRoutine godoc
//
// @Summary		Delete Routine API
// @Description	Delete a routine by ID.
// @Accept		json
// @Produce		json
// @Tags 		데일리 루틴
// @Param		id	path	string	true	"Routine ID"
// @Success		204
// @Failure		409	{object}	errorcode.UseCaseError
// @Router		/v1/routines/{id} [delete]
func (svc RoutineUseCase) DeleteRoutine(c *fiber.Ctx, routineId uuid.UUID) error {
	err := svc.routineOutputPort.DeleteRoutine(routineId)
	if err != nil {
		return err
	}
	return nil
}

// WriteRoutineRecord godoc
//
// @Summary		Write Routine Record API
// @Description	Write a record for a routine.
// @Accept		json
// @Produce		json
// @Tags 		데일리 루틴
// @Param		id	path	string	true	"Routine ID"
// @Param		messages.WriteRoutineRecordRequest	body	messages.WriteRoutineRecordRequest	true	"Write Routine Record Request"
// @Success		200	{object}	messages.RoutineRecordResponse
// @Failure		409	{object}	errorcode.UseCaseError
// @Router		/v1/routines/{id}/record [post]
func (svc RoutineUseCase) WriteRoutineRecord(c *fiber.Ctx, req messages.WriteRoutineRecordRequest) (*messages.RoutineRecordResponse, error) {
	routine, err := svc.routineOutputPort.GetRoutineById(req.RoutineId)
	if err != nil {
		return nil, errorcode.ErrRoutineNotFound
	}

	newRoutineRecord := &domains.RoutineRecord{
		RecordContent: req.RecordContent,
		IsRecord:      true,
	}

	// 루틴 세트 조회
	routineSet, err := svc.routineOutputPort.GetRoutineSetByRoutineIdAndDate(req.UserId, routine.ID, req.Year, req.Month, req.Day)
	if err != nil { // 루틴 세트가 없으면 루틴 세트를 생성
		routines, err := svc.routineOutputPort.ListRoutinesByUserId(req.UserId)
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

// DeleteRoutineRecord godoc
//
// @Summary		Delete Routine Record API
// @Description	Delete a routine record by ID.
// @Accept		json
// @Produce		json
// @Tags 		데일리 루틴
// @Param		id	path	string	true	"Routine Record ID"
// @Param		messages.DeleteRoutineRecordRequest	body	messages.DeleteRoutineRecordRequest	true	"Delete Routine Record Request"
// @Success		200
// @Failure		409	{object}	errorcode.UseCaseError
// @Router		/v1/routines/record/{id} [delete]
func (svc RoutineUseCase) DeleteRoutineRecord(c *fiber.Ctx, req messages.DeleteRoutineRecordRequest) error {
	record, err := svc.routineOutputPort.GetRoutineRecordById(req.RoutineRecordId)
	if err != nil {
		return errorcode.ErrRoutineNotFound
	}

	if record.RoutineSet.UserId != req.UserId {
		return errorcode.ErrRoutineRecordNotFound
	}

	return svc.routineOutputPort.DeleteRoutineRecord(req.RoutineRecordId)
}
