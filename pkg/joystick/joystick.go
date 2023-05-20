package joystick

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"os"
)

type Joystick struct {
	path        string
	file        *os.File
	axisCount   int
	buttonCount int
}

func NewJoystick(path string) *Joystick {
	return &Joystick{
		path: path,
	}
}

func (j *Joystick) Open() error {
	if j.file != nil {
		_ = j.file.Close()
	}

	f, err := os.OpenFile(j.path, os.O_RDONLY, 0666)
	if err != nil {
		return fmt.Errorf("failed to open joystick %s: %w", j.path, err)
	}
	j.file = f

	j.buttonCount, err = unix.IoctlGetInt(int(j.file.Fd()), JSIOCGBUTTONS)
	if err != nil {
		return fmt.Errorf("failed to get button count of %s: %w", j.path, err)
	}
	j.axisCount, err = unix.IoctlGetInt(int(j.file.Fd()), JSIOCGAXES)
	if err != nil {
		return fmt.Errorf("failed to get axis count of %s: %w", j.path, err)
	}

	log.Printf("opened joystick with %d axis and %d buttons\n", j.axisCount, j.buttonCount)

	return nil
}

func (j *Joystick) Run() <-chan js_event {
	if j.file == nil {
		return nil
	}

	ch := make(chan js_event)

	go func() {
		for {
			b := make([]byte, 8)
			i, err := j.file.Read(b)
			if err != nil {
				break
			}
			if i != 8 {
				break
			}
			data := bytes.NewReader(b)
			var ev js_event
			err = binary.Read(data, binary.LittleEndian, &ev)
			ch <- ev
		}
		_ = j.file.Close()
		close(ch)
	}()

	return ch
}
