package photo

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func (s *Service) HandleUploadPhotos(files []*multipart.FileHeader) (interface{}, error) {
	for _, file := range files {
		log.Println(file.Filename)

		err := s.saveUploadFile(file)
		if err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"message": "Upload file successfully",
	}, nil
}

func (s *Service) GetByID(id string) (interface{}, error) {
	buffer, err := os.ReadFile("uploads/" + id)

	if err != nil {
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString(buffer)

	return map[string]interface{}{
		"image": encoded,
	}, nil
}

func (s *Service) HandlerGetUsers() (interface{}, error) {
	var data, err = s.repo.List()

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Service) HandlerGetUser(id string) (*Photo, error) {
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

func (s *Service) HandlerCreateUser(body *CreateUserReq) interface{} {
	// result := s.repo.InsertOne(body.Name, body.Description)

	return map[string]string{
		"id": "!",
	}
}

func (s *Service) HandlerUpdateUser(body *CreateUserReq) interface{} {
	// result := s.repo.InsertOne(body.Name, body.Description)

	return map[string]string{
		"id": "!",
	}
}

func (s *Service) HandlerGetProfile(id string) (interface{}, error) {
	var data struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	params := &QueryParams{
		TableName: "users",
		Columns:   []string{"id", "name", "email"},
		Where:     "id = $1",
		Args:      []interface{}{id},
	}

	if err := s.repo.DetailByConditions(&data, params); err != nil {
		return nil, err
	}

	return &data, nil
}

func (s *Service) HandlerUpdateProfile(id string, body interface{}) (interface{}, error) {
	var data struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	params := &QueryParams{
		TableName: "users",
		Columns:   []string{"id", "name", "email"},
		Where:     "id = $1",
		Args:      []interface{}{id},
	}

	if err := s.repo.DetailByConditions(&data, params); err != nil {
		return nil, err
	}

	return &data, nil
}

func (s *Service) saveUploadFile(file *multipart.FileHeader) error {
	// Open the uploaded file
	fileReader, err := file.Open()
	if err != nil {
		fmt.Println(err.Error(), 123)
		return err
	}
	// defer fileReader.Close()

	// Read the file contents into a buffer
	buffer := bytes.NewBuffer(nil)

	_, err = io.Copy(buffer, fileReader)
	if err != nil {
		return err
	}

	// Save the buffer content to a file
	err = os.WriteFile("uploads/"+uuid.NewString(), buffer.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}
