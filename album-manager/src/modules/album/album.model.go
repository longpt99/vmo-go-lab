package album

type CreateAlbumReq struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type UpdateAlbumReq struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}
