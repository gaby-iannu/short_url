package repository

import "short_url/short/model"


var urls map[string]model.Url= make(map[string]model.Url)

type Repository interface {
	InsertIfNotExists(tiny string, url model.Url) bool
	Read(tinyUrl string) model.Url
}

type repository struct {}

func New() Repository {
	return &repository{}
}


// true si lo inserta 
// false si ya existe y no lo inserta
func (r *repository) InsertIfNotExists(tiny string, url model.Url) bool {
	
	if _,ok := urls[tiny]; !ok {
		urls[tiny] = url
		return true
	}

	return false
}

func (r *repository) Read(tinyUrl string) model.Url {
	if v,ok := urls[tinyUrl]; ok {
		return v
	}

	return model.Url{}
}
