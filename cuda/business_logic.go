package cuda

import (
	"errors"
)

func ApplyBusinessRules(data DataSource) (DataSource, error) {
	if data.Description == "" {
		return DataSource{}, errors.New("Description field is empty")
	}
	return data, nil
}
