package model

type Url struct {
	Long string `json:"long_url"` 
	User string `json:"user_id"`
}


func New(long, user string) Url {
	return Url{
		Long: long,
		User: user,
	}
}
