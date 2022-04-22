package voo

const version = "0.0.1"

type Voo struct {
	AppName string
	Debug   bool
	Version string
}

func (v *Voo) New(rootPath string) error {
	pathConfig := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"migrations", "handlers", "logs", "data", "tmp", "views", "public", "tmp", "middlewares"},
	}
	err := v.Init(pathConfig)
	if err != nil {

		return err
	}
	return nil
}

func (v *Voo) Init(p initPaths) error {
	root := p.rootPath
	for _, path := range p.folderNames {
		//! if folder does not exits, create it
		err := v.CreateDirIfNotExits(root + "/" + path)
		if err != nil {
			return err
		}
	}

	return nil
}
