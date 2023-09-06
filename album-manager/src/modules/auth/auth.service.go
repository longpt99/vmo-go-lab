package auth

import (
	mC "album-manager/src/common/models"
	"album-manager/src/common/repository"
	"album-manager/src/configs/database"
	"album-manager/src/errors"
	"album-manager/src/models"
	"album-manager/src/modules/notification"
	"album-manager/src/modules/user"
	"album-manager/src/utils"
	"album-manager/src/utils/object"
	"album-manager/src/utils/strings"
	t "album-manager/src/utils/token"
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"
)

type Service struct {
	// repo Repository
	userRepo user.Repository
}

func (s *Service) handleSignIn(body *models.LoginUserReq) (interface{}, error) {
	op := errors.Op("auth.service.handleSignIn")

	params := &repository.QueryParams{
		Columns: []string{"id", "password", "status", "email"},
		Where:   "email ILIKE ? OR username ILIKE ?",
		Args:    []interface{}{body.Identifier, body.Identifier},
	}

	doc, err := s.userRepo.DetailByConditions(params)
	if err != nil {
		return nil, err
	}

	if reflect.ValueOf(doc).IsZero() {
		return nil, errors.E(op, http.StatusBadRequest, "account not found")
	}

	if doc.Status == mC.INACTIVE {
		OTP := strings.RandomNumStr(6)
		notiService := notification.Service{
			Notifier: notification.CreateNotifier("email"),
		}

		go notiService.SendNotification(fmt.Sprintf("Your OTP code is %s", OTP), []string{
			doc.Email,
		})

		database.RedisClient.SetEx(context.Background(), fmt.Sprintf("caches:users:%s:active_account", doc.ID), OTP, time.Minute*5)

		return nil, errors.E(op, http.StatusBadRequest, "your account hasn't been activated. Please check your email!")
	}

	match := utils.CompareHashPassword(body.Password, doc.Password)
	if !match {
		return nil, errors.E(op, http.StatusBadRequest, "wrong password!")
	}

	return map[string]string{"access_token": t.SignToken(doc.ID)}, nil
}

func (s *Service) handleActiveAccount(body *models.ActiveAccountReq) (interface{}, error) {
	op := errors.Op("auth.service.handleActiveAccount")

	params := &repository.QueryParams{
		Columns: []string{"id", "status"},
		Where:   "email ILIKE ?",
		Args:    []interface{}{fmt.Sprintf(`%s`, body.Email)},
	}

	doc, err := s.userRepo.DetailByConditions(params)
	if err != nil {
		return nil, err
	}

	if reflect.ValueOf(doc).IsZero() {
		return nil, errors.E(op, http.StatusBadRequest, "account not found")
	}

	if doc.Status != mC.INACTIVE {
		return nil, errors.E(op, http.StatusBadRequest, "your account has been activated.")
	}

	var (
		OTP      string
		redisKey = fmt.Sprintf("caches:users:%s:active_account", doc.ID)
	)

	err = database.RedisClient.Get(context.Background(), redisKey).Scan(&OTP)
	if err != nil {
		return nil, errors.E(op, http.StatusBadRequest, "OTP is invalid.")
	}

	if OTP != body.OTP {
		return nil, errors.E(op, http.StatusBadRequest, "OTP is not match.")
	}

	updateParams := &repository.UpdateParams{
		TableName: "users",
		Where:     "id = ?",
		Args:      []interface{}{doc.ID},
		Data: map[string]interface{}{
			"status": mC.ACTIVE,
		},
	}

	_ = s.userRepo.UpdateByConditions(updateParams)

	database.RedisClient.Del(context.Background(), redisKey)

	return nil, nil
}

func (s *Service) handleSignUp(body *models.SignUpUserReq) (interface{}, error) {
	op := errors.Op("auth.service.handleSignUp")

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

	// TODO: Send link to active account
	// notiService := notification.Service{
	// 	Notifier: notification.CreateNotifier("email"),
	// }
	// notiService.SendNotification("OTP: ", []string{
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
