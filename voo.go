package voo

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/VooDooStack/Voo/render"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

const version = "0.0.1"

type Voo struct {
	AppName  string
	Debug    bool
	Version  string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	RootPath string
	Routes   *chi.Mux
	Render   *render.Render
	config   config
}

type config struct {
	port     string
	renderer string
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

	err = v.checkDotEnv(rootPath)
	if err != nil {
		return err
	}

	//! read .env file
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	//! create loggers
	infoLog, errorLog := v.startLoggers()
	v.InfoLog = infoLog
	v.ErrorLog = errorLog
	v.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	v.Version = version
	v.RootPath = rootPath
	v.Routes = v.routes().(*chi.Mux)

	v.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
	}

	v.Render = v.createRenderer(v)

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

//! listen and serve starts server
func (v *Voo) ListenAndServe() {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog:     v.ErrorLog,
		Handler:      v.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	v.InfoLog.Printf("listening on port %s", os.Getenv("PORT"))
	err := srv.ListenAndServe()
	if err != nil {
		v.ErrorLog.Fatal(err)
	}
}

func (v *Voo) checkDotEnv(path string) error {
	err := v.CreateFileIfNotExists(fmt.Sprintf("%s/.env", path))
	if err != nil {
		return err
	}
	return nil
}

func (v *Voo) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t ", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}

func (v *Voo) createRenderer(voo *Voo) *render.Render {
	myRenderer := render.Render{
		Renderer: voo.config.renderer,
		RootPath: voo.RootPath,
		Port:     voo.config.port,
	}

	return &myRenderer
}
