package helper

import (
	"errors"
	"fmt"
	"os"

	"github.com/rahmatrdn/go-skeleton/config"
	apperr "github.com/rahmatrdn/go-skeleton/error"
)

func CreateDirectory(id_merchant, filename string) string {
	filename = DateFilename() + "-" + filename
	storagePath := config.StorageDirectory + id_merchant

	if _, err := os.Stat(storagePath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(storagePath, 0755)
		if err != nil {
			apperr.PanicLogging(err)
		}
	}

	return fmt.Sprintf(storagePath+"/%s", filename)
}
