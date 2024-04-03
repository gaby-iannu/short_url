package model

type Url struct {
	Long string `json:"long_url" binding:"required"` 
	User string `json:"user_id" binding:"required"`
}


func New(long, user string) Url {
	return Url{
		Long: long,
		User: user,
	}
}
