package model

type Url struct {
	Long string
	User string
}


func New(long, user string) Url {
	return Url{
		Long: long,
		User: user,
	}
}
