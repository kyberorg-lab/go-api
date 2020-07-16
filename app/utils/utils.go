package utils

type ErrJson struct {
	Err   string `json:"err"`
	Error error  `json:"error" binding: "omitempty"`
}

func ErrorJson(err string) ErrJson {
	return ErrJson{err, nil}
}

func ErrorJsonWithError(message string, err error) ErrJson {
	return ErrJson{
		Err:   message,
		Error: err,
	}
}
