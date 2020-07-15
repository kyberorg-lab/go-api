package utils

type ErrJson struct {
	Err string `json:"err"`
}

func ErrorJson(err string) ErrJson {
	return ErrJson{err}
}
