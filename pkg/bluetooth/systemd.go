package bluetooth

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

const (
	dropInDir  = "/etc/systemd/system/bluetooth.service.d"
	dropInPath = dropInDir + "/deckjoy.conf"

	dropInContent = `[Service]
ExecStart=
ExecStart=/usr/lib/bluetooth/bluetoothd -P input
`
)

func EnsureBluetoothdInputDisabled() error {
	existing, err := os.ReadFile(dropInPath)
	if err == nil && bytes.Equal(existing, []byte(dropInContent)) {
		return nil
	}
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read %s: %w", dropInPath, err)
	}

	if err := os.MkdirAll(filepath.Dir(dropInPath), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", dropInDir, err)
	}
	if err := os.WriteFile(dropInPath, []byte(dropInContent), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", dropInPath, err)
	}
	log.Infof("installed bluetoothd drop-in %s", dropInPath)

	if out, err := exec.Command("systemctl", "daemon-reload").CombinedOutput(); err != nil {
		return fmt.Errorf("systemctl daemon-reload: %w: %s", err, out)
	}
	if out, err := exec.Command("systemctl", "restart", "bluetooth.service").CombinedOutput(); err != nil {
		return fmt.Errorf("systemctl restart bluetooth: %w: %s", err, out)
	}
	log.Infof("restarted bluetooth.service")
	return nil
}
