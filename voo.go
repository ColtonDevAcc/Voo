package voo

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/VooDooStack/Voo/render"
	"github.com/VooDooStack/Voo/session"
	"github.com/alexedwards/scs/v2"
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
	Session  *scs.SessionManager
	DB       Database
	JetViews *jet.Set
	config   config
}

type config struct {
	port        string
	renderer    string
	cookie      cookieConfig
	sessionType string
	database    databaseConfig
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

	//! connect to database
	if os.Getenv("DATABASE_TYPE") != "" {
		db, err := v.openDB(os.Getenv("DATABASE_TYPE"), v.BuildDSN())
		if err != nil {
			errorLog.Println(err)
			os.Exit(1)
		}

		v.DB = Database{
			DatabseType: os.Getenv("DATABASE_TYPE"),
			Pool:        db,
		}
	}

	v.InfoLog = infoLog
	v.ErrorLog = errorLog
	v.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	v.Version = version
	v.RootPath = rootPath
	v.Routes = v.routes().(*chi.Mux)

	v.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
		cookie: cookieConfig{
			name:     os.Getenv("COOKIE_NAME"),
			lifetime: os.Getenv("COOKIE_LIFETIME"),
			persist:  os.Getenv("COOKIE_PERSIST"),
			secure:   os.Getenv("COOKIE_SECURE"),
			domain:   os.Getenv("COOKIE_DOMAIN"),
		},
		sessionType: os.Getenv("SESSION_TYPE"),
		database: databaseConfig{
			databse: os.Getenv("DATABASE_TYPE"),
			dsn:     v.BuildDSN(),
		},
	}

	//! Create a session
	sess := session.Session{
		CookieLifetime: v.config.cookie.lifetime,
		CookiePersist:  v.config.cookie.persist,
		CookieName:     v.config.cookie.name,
		SessionType:    v.config.sessionType,
		CookieDomain:   v.config.cookie.domain,
	}

	v.Session = sess.InitSession()

	var views = jet.NewSet(
		jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rootPath)),
		jet.InDevelopmentMode(),
	)

	v.JetViews = views

	v.createRenderer()

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

func (v *Voo) createRenderer() {
	myRenderer := render.Render{
		Renderer: v.config.renderer,
		RootPath: v.RootPath,
		Port:     v.config.port,
		JetViews: v.JetViews,
	}
	v.Render = &myRenderer
}

func (v *Voo) BuildDSN() string {
	var dsn string

	switch os.Getenv("DATABASE_TYPE") {
	case "postgres", "postgresql":
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5",
			os.Getenv("DATABASE_HOST"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_NAME"),
			os.Getenv("DATABASE_SSL_MODE"))

		// we check to see if a database passsword has been supplied, since including "password=" with nothing
		// after it sometimes causes postgres to fail to allow a connection.
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("%s password=%s", dsn, os.Getenv("DATABASE_PASS"))
		}

	default:

	}

	return dsn
}
