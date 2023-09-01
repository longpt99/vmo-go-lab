package auth

import (
	"album-manager/src/common/repository"
	"album-manager/src/errors"
	"album-manager/src/models"
	"album-manager/src/modules/user"
	"album-manager/src/utils"
	"album-manager/src/utils/object"
	t "album-manager/src/utils/token"
	"net/http"
	"reflect"
)

type Service struct {
	// repo Repository
	userRepo user.Repository
}

func (s *Service) handleSignIn(body *models.LoginUserReq) (interface{}, error) {
	op := errors.Op("auth.service.handleSignIn")

	// var data struct {
	// 	ID       string `json:"id"`
	// 	Password string `json:"password"`
	// }

	params := &repository.QueryParams{
		Columns: []string{"id", "password"},
		Where:   "email ILIKE ? OR username ILIKE ?",
		Args:    []interface{}{body.Identifier, body.Identifier},
	}

	result, err := s.userRepo.DetailByConditions(params)
	if err != nil {
		return nil, err
	}

	if reflect.ValueOf(result).IsZero() {
		return nil, errors.E(op, http.StatusBadRequest, "account not found")
	}

	match := utils.CompareHashPassword(body.Password, result.Password)
	if !match {
		return nil, errors.E(op, http.StatusBadRequest, "wrong password!")
	}

	return map[string]string{"access_token": t.SignToken(result.ID)}, nil
}

func (s *Service) handleSignUp(body *models.SignUpUserReq) (interface{}, error) {
	op := errors.Op("auth.service.handleSignUp")

	// var data struct {
	// 	ID string `json:"id"`
	// }

	params := &repository.QueryParams{
		TableName: "users",
		Columns:   []string{"id"},
		Where:     "email ILIKE ?",
		Args:      []interface{}{body.Email},
	}

	result, err := s.userRepo.DetailByConditions(params)
	if err != nil {
		return nil, err
	}

	if !reflect.ValueOf(result).IsZero() {
		return nil, errors.E(op, http.StatusBadRequest, "account have exists")
	}

	createUserParams := &models.User{}

	err = object.MergeStructIntoModel(createUserParams, body)
	if err != nil {
		return nil, err
	}

	createUserParams.Password = utils.HashPassword(body.Password)

	id, err := s.userRepo.InsertOne(createUserParams)
	if err != nil {
		return nil, err
	}

	//TODO: Send link to active account
	// notiService := notification.Service{
	// 	notification.CreateNotifier("email"),
	// }
	// notiService.SendNotification("Hello sờ lyly!", []string{
	// 	body.Email,
	// })

	return map[string]interface{}{
		"id": id,
	}, nil
}

func (s *Service) handleResetPassword(body *models.ResetPasswordReq) (interface{}, error) {
	op := errors.Op("auth.service.handleSignUp")

	params := &repository.QueryParams{
		TableName: "users",
		Columns:   []string{"id"},
		Where:     "email ILIKE ? OR username ILIKE ?",
		Args:      []interface{}{body.Identifier, body.Identifier},
	}

	count, err := s.userRepo.CountByConditions(params)
	if err != nil {
		return nil, err
	}

	if *count == 0 {
		return nil, errors.E(op, http.StatusBadRequest, "account is not exists")
	}

	//TODO Send new password to mail

	return map[string]bool{
		"isSucceed": true,
	}, nil
}

func (s *Service) handleChangePassword(body *models.ChangePasswordReq, id string) (interface{}, error) {
	op := errors.Op("auth.service.handleChangePassword")

	// var data struct {
	// 	ID       string `json:"id"`
	// 	Password string `json:"password"`
	// }

	params := &repository.QueryParams{
		TableName: "users",
		Columns:   []string{"id"},
		Where:     "id = ?",
		Args:      []interface{}{id},
	}

	result, err := s.userRepo.DetailByConditions(params)
	if err != nil {
		return nil, err
	}

	match := utils.CompareHashPassword(body.OldPassword, result.Password)
	if !match {
		return nil, errors.E(op, http.StatusBadRequest, "old password not match")
	}

	updateParams := &repository.UpdateParams{
		TableName: "users",
		Where:     "id = ?",
		Args:      []interface{}{id},
		Data: map[string]interface{}{
			"password": utils.HashPassword(body.Password),
		},
	}

	err = s.userRepo.UpdateByConditions(updateParams)
	if err != nil {
		return nil, err
	}

	//TODO Send new password to mail
	return map[string]interface{}{
		"data": result,
	}, nil
}
