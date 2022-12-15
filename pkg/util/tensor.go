package util

import (
	"errors"
	"fmt"
	"os"
)

func TensorBytesToFile(sceneID string, tensor []byte) (string, error) {
	tensorPath := fmt.Sprintf("/tmp/%s", sceneID)

	if err := os.WriteFile(tensorPath, tensor, 0666); err != nil {
		return "", err
	}

	return tensorPath, nil
}

func TensorCachedPath(sceneID string) string {
	tensorPath := fmt.Sprintf("/tmp/%s", sceneID)

	if _, err := os.Stat(tensorPath); errors.Is(err, os.ErrNotExist) {
		return ""
	}

	return tensorPath
}
