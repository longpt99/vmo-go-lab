package album

import (
	mC "album-manager/src/common/models"
	"album-manager/src/common/repository"
	"album-manager/src/errors"
	"album-manager/src/models"
	"album-manager/src/modules/user"
	"album-manager/src/utils/object"
	"net/http"
)

type Service struct {
	repo  Repository
	userR user.Repository
}

func (s *Service) List(userID string, query mC.QueryStringParams) (interface{}, error) {
	var data, err = s.repo.List(&repository.FindParams{
		Where:             "m2m.user_id = ?",
		Args:              []interface{}{userID},
		Select:            []string{"albums.id", "albums.name", "albums.description"},
		Joins:             []string{"LEFT JOIN user_albums AS m2m ON m2m.album_id = albums.id", "LEFT JOIN users ON users.id = m2m.user_id"},
		QueryStringParams: query,
	})

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
		return nil, errors.Str("User Not Found")
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
	var owner = models.User{ID: userID}

	params := &models.Album{
		Users:   []models.User{owner},
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

func (s *Service) updateByID(userID, id string, body *UpdateAlbumReq) (interface{}, error) {
	count, err := s.repo.CountByConditions(&repository.QueryParams{
		Where: "owner_id = ? AND id = ?",
		Args:  []interface{}{userID, id},
	})
	if err != nil {
		return nil, err
	}

	if *count == 0 {
		return nil, errors.E(errors.Op("updateByID"), http.StatusBadRequest, "album not found")
	}

	params := &models.Album{}

	err = object.MergeStructIntoModel(params, body)
	if err != nil {
		return nil, err
	}

	err = s.repo.UpdateByConditions(&repository.UpdateParams{
		Where: "owner_id = ? AND id = ?",
		Args:  []interface{}{userID, id},
		Data:  params,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *Service) deleteByID(userID, id string) (interface{}, error) {
	count, err := s.repo.CountByConditions(&repository.QueryParams{
		Where: "owner_id = ? AND id = ?",
		Args:  []interface{}{userID, id},
	})
	if err != nil {
		return nil, err
	}

	if *count == 0 {
		return nil, errors.E(errors.Op("deleteByID"), http.StatusBadRequest, "album not found")
	}

	err = s.repo.Delete(id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
