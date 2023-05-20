package usbgadget

import "path"

type FunctionECM struct {
	*FunctionGeneric
}

type FunctionECMAttributes struct {
	HostAddr string
	DevAddr  string
}

func (g *Gadget) CreateFunctionECM(name string) (*FunctionECM, error) {
	gf := g.GetFunction(name)
	f := &FunctionECM{
		gf,
	}

	if err := mkDir(f.path); err != nil {
		return nil, err
	}

	return f, nil
}

func (f *FunctionECM) SetAttributes(attrs FunctionECMAttributes) error {
	if err := writeString(path.Join(f.path, "host_addr"), attrs.HostAddr); err != nil {
		return err
	}
	if err := writeString(path.Join(f.path, "dev_addr"), attrs.DevAddr); err != nil {
		return err
	}
	return nil
}
