package user

type UpdateUserProfileReq struct {
	Name string `json:"name"`
	// Email string `json:"email" validate:"email"`
	// DOB   string `json:"dob" validate:"date_string"`
}
