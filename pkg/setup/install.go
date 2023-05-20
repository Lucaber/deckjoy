package setup

import (
	"fmt"
	"github.com/lucaber/deckjoy/pkg/util"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

func Install() error {
	if err := installControllerConfig(); err != nil {
		log.WithError(err).Error("failed to install controller config")
	}
	return nil
}

const controllerConfigFileName = "controller_neptune.vdf"
const steamControllerConfigsPath = ".local/share/Steam/steamapps/common/Steam Controller Configs/"
const steamControllerConfigSubPath = "config/deckjoy/"

func installControllerConfig() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	sourcePath := path.Join(cwd, controllerConfigFileName)
	if _, err := os.Stat(sourcePath); err != nil {
		return fmt.Errorf("file not found %s %w", sourcePath, err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	userids, err := os.ReadDir(path.Join(home, steamControllerConfigsPath))
	if err != nil {
		return err
	}
	for _, userid := range userids {
		path := path.Join(home, steamControllerConfigsPath, userid.Name(), steamControllerConfigSubPath, controllerConfigFileName)
		if _, err := os.Stat(path); err == nil {
			// already exists

			//continue
		}
		err = util.CopyFile(sourcePath, path)
		if err != nil {
			return err
		}
	}

	return nil
}
