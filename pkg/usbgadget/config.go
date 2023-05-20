package usbgadget

import (
	"os"
	"path"
)

type Config struct {
	path string
}

type ConfigAttributes struct {
	MaxPower uint16
}

type ConfigStrings struct {
	Configuration string
}

func (g *Gadget) CreateConfig(name string) (*Config, error) {
	c := g.GetConfig(name)

	if err := mkDir(c.path); err != nil {
		return nil, err
	}

	return c, nil
}

func (g *Gadget) GetConfig(name string) *Config {
	return &Config{
		path: path.Join(g.path, "configs", name),
	}
}

func (g *Gadget) GetConfigs() ([]*Config, error) {
	entries, err := listDir(path.Join(g.path, "configs"))
	if err != nil {
		return nil, err
	}
	configs := make([]*Config, len(entries))
	for i, e := range entries {
		configs[i] = g.GetConfig(e)
	}

	return configs, nil
}

func (c *Config) SetAttributes(attrs ConfigAttributes) error {
	if err := writeHex16(path.Join(c.path, "MaxPower"), attrs.MaxPower); err != nil {
		return err
	}
	return nil
}
func (c *Config) SetStrings(strs ConfigStrings) error {
	stringsPath := path.Join(c.path, "strings", "0x409")
	_ = mkDir(stringsPath)

	if err := writeString(path.Join(stringsPath, "configuration"), strs.Configuration); err != nil {
		return err
	}

	return nil
}

func (c *Config) AddFunction(f *FunctionGeneric) error {
	return os.Symlink(f.path, path.Join(c.path, f.name))
}

func (c *Config) RemoveFunction(f *FunctionGeneric) error {
	return os.Remove(path.Join(c.path, f.name))
}

func (c *Config) DeleteStrings() error {
	stringsPath := path.Join(c.path, "strings", "0x409")
	return os.Remove(stringsPath)
}

func (c *Config) RemoveAllFunctions() error {
	entries, err := os.ReadDir(c.path)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.Type() == os.ModeSymlink {
			err := os.Remove(path.Join(c.path, e.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Config) Destroy(skipErrors bool) error {
	if err := c.RemoveAllFunctions(); err != nil {
		if !skipErrors {
			return err
		}
	}

	if err := c.DeleteStrings(); err != nil {
		if !skipErrors {
			return err
		}
	}

	if err := os.Remove(c.path); err != nil {
		if !skipErrors {
			return err
		}
	}
	return nil
}
