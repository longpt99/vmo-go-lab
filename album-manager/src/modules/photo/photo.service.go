package photo

import (
	"bytes"
	"encoding/base64"
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
