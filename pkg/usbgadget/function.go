package usbgadget

import (
	"os"
	"path"
)

type FunctionGeneric struct {
	path string
	name string
}

func (g *Gadget) GetFunction(name string) *FunctionGeneric {
	return &FunctionGeneric{
		name: name,
		path: path.Join(g.path, "functions", name),
	}
}

func (g *Gadget) GetFunctions() ([]*FunctionGeneric, error) {
	entries, err := listDir(path.Join(g.path, "functions"))
	if err != nil {
		return nil, err
	}
	functions := make([]*FunctionGeneric, len(entries))
	for i, e := range entries {
		functions[i] = g.GetFunction(e)
	}

	return functions, nil
}

func (f *FunctionGeneric) Destroy(skipErrors bool) error {
	if err := os.Remove(f.path); err != nil {
		if !skipErrors {
			return err
		}
	}
	return nil
}
