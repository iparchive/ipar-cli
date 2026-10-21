package ipar

import (
	"reflect"

	"github.com/chigopher/pathlib"
	"github.com/go-viper/mapstructure/v2"
)

func StringToPath() mapstructure.DecodeHookFuncType {
	return func(
		from reflect.Type,
		to reflect.Type,
		data any,
	) (any, error) {
		if from != reflect.TypeFor[string]() {
			return data, nil
		}
		if to != reflect.TypeFor[pathlib.Path]() {
			return data, nil
		}
		return pathlib.NewPath(data.(string)), nil
	}
}
