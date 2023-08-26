package album

import (
	"album-manager/src/models"
	"album-manager/src/modules/user"
	"album-manager/src/utils/object"
	"errors"
)

type Service struct {
	repo  Repository
	userR user.Repository
}

func (s *Service) HandlerGetUsers() (interface{}, error) {
	var data, err = s.repo.List()

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Service) HandlerGetUser(id string) (*models.Album, error) {
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

func (s *Service) create(userID string, body *CreateAlbumReq) (interface{}, error) {
	var owner = &models.User{ID: userID}

	params := &models.Album{
		Users:   []*models.User{owner},
		OwnerID: &userID,
	}

	err := object.MergeStructIntoModel(params, body)
	if err != nil {
		return nil, err
	}

	id, err := s.repo.InsertOne(params)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id": id,
	}, nil
}

func (s *Service) deleteByID(userID, id string) (interface{}, error) {

	// count := s.repo.

	// err := s.repo.Delete(id)
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}
