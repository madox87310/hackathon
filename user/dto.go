package user

type SignUpRequest struct {
	DisplayName string `json:"display_name" binding:"required,min=1,max=32"`
	PhoneNumber string `json:"phone_number" binding:"required,e164"`
	Password    string `json:"password" binding:"required,min=8,max=72"`
}

type SignUpResponse struct {
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SignInRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required,e164"`
	Password    string `json:"password" binding:"required,min=8,max=72"`
}

type SignInResponse struct {
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
