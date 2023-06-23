package file

import (
	"argus/pkg/models"
	"argus/pkg/storage/schema"
	"encoding/gob"
	"fmt"
	"os"
)

func init() {
	schema.Register("file", &FileStorage{})
}

type FileStorage struct {
	fileName string
}

func (f *FileStorage) Save(config *models.Configuration) error {
	dataFile, err := os.Create(f.fileName)
	if err != nil {
		return fmt.Errorf("Could not create file %v:%w", f.fileName, err)
	}
	dataEncoder := gob.NewEncoder(dataFile)
	return dataEncoder.Encode(config)
}

func (f *FileStorage) Configure(config map[string]interface{}) error {
	var ok bool
	f.fileName, ok = config["file"].(string)
	if !ok {
		return fmt.Errorf("could not configure file path. is this name valid?")
	}
	return nil
}

func (f *FileStorage) Load() (*models.Configuration, error) {
	dataFile, err := os.Open(f.fileName)
	if err != nil {
		return nil, fmt.Errorf("Could not open file %v:%w", f.fileName, err)
	}
	dataDecoder := gob.NewDecoder(dataFile)
	var config models.Configuration
	err = dataDecoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("Could not decode file %v:%w", f.fileName, err)
	}
	return &config, nil
}
