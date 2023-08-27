package album

type CreateAlbumReq struct {
	Name        string `json:"name" validate:"required,ascii"`
	Description string `json:"description" validate:"omitempty,ascii"`
}

type UpdateAlbumReq struct {
	Name        string `json:"name" validate:"omitempty,alphanumspace"`
	Description string `json:"description" validate:"omitempty,ascii"`
}
