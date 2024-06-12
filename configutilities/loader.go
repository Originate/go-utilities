package configutilities

import (
	"errors"
	"flag"
	"io/fs"
	"log/slog"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

const defaultFileLocation = "./config/config.yml"

func Load(cfg interface{}) error {
	filePath := flag.String("config", defaultFileLocation, "The custom config file for the app to use")
	flag.Parse()

	if reflect.ValueOf(cfg).Kind() != reflect.Ptr {
		return errors.New("config can't be loaded on a non pointer struct")
	}

	if err := cleanenv.ReadConfig(*filePath, cfg); err != nil {
		if *filePath == defaultFileLocation && errors.Is(err, fs.ErrNotExist) {
			slog.Debug("Default config file is missing, reading from env")

			if err := cleanenv.ReadEnv(cfg); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return registerValidationAliases(validator.New(validator.WithRequiredStructEnabled())).Struct(cfg)
}
