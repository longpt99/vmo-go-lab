package user

import (
	"album-manager/src/common/repository"
	"album-manager/src/models"
	"errors"
)

type Service struct {
	repo Repository
}

func (s *Service) HandlerGetUsers() (interface{}, error) {
	var data, err = s.repo.List(&repository.FindParams{})

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Service) HandlerGetUser(id string) (*models.User, error) {
	var data, err = s.repo.DetailByID(id)

	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, errors.New("User Not Found")
	}

	return data, nil
}

func (s *Service) HandlerDeleteUser(id string) (interface{}, error) {
	err := s.repo.Delete(id)

	if err != nil {
		return nil, err
	}

	return map[string]bool{
		"is_succeed": true,
	}, nil
}

func (s *Service) HandlerCreateUser(body *models.User) interface{} {
	// result := s.repo.InsertOne(body.Name, body.Description)

	return map[string]string{
		"id": "!",
	}
}

func (s *Service) HandlerUpdateUser(body *models.CreateUserReq) interface{} {
	// result := s.repo.InsertOne(body.Name, body.Description)

	return map[string]string{
		"id": "!",
	}
}

func (s *Service) HandlerGetProfile(id string) (interface{}, error) {
	// var data struct {
	// 	ID       string `json:"id"`
	// 	Email    string `json:"email"`
	// 	Name     string `json:"name"`
	// 	Username string `json:"username"`
	// }

	params := &repository.QueryParams{
		TableName: "users",
		Columns:   []string{"id", "name", "email", "username"},
		Where:     "id = $1",
		Args:      []interface{}{id},
	}

	result, err := s.repo.DetailByConditions(params)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *Service) HandlerUpdateProfile(id string, body UpdateUserProfileReq) (interface{}, error) {
	params := &repository.UpdateParams{
		TableName: "users",
		Where:     "id = ?",
		Args:      []interface{}{id},
		Data:      body,
	}

	if err := s.repo.UpdateByConditions(params); err != nil {
		return nil, err
	}

	return nil, nil
}
