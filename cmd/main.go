package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"

	"webpage/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

// type config struct{
// 	Addr string
// 	Staticdir string
// }

type application struct{
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {

	// addr := flag.String("addr",":4000","HTTP network address")
	// flag.Parse()

	
	//not good we have to impliclity decleare in environment varaible
	// addr := os.Getenv("SNIPPETBOX_ADDR")


	// cfg := new(config)
	// flag.StringVar(&cfg.Addr,"addr",":4000","HTTP new address")
	// flag.StringVar(&cfg.Staticdir,"static-dir","../ui/static","path to new directory")
	// flag.Parse()

	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn","web:Admin@123@/snippetbox?parseTime=true","MySQL data source name")
	flag.Parse()


	infoLog := log.New(os.Stdout,"INFO\t",log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr,"ERROR\t",log.Ldate|log.Ltime|log.Lshortfile)

	db,err := openDB(*dsn)
	if err != nil{
		errorLog.Fatal(err)
	}

	defer db.Close()


	templateCache, err := newTemplateCache("../ui/html")
	if err != nil{
		errorLog.Fatal(err)
	}


	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}



	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)


	fileServer := http.FileServer(http.Dir("../ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: mux,
	}
	

	infoLog.Printf("Starting server on %s",*addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string)(*sql.DB,error){
	db,err := sql.Open("mysql",dsn)
	if err != nil{
		
		return nil,err
	}
	if err = db.Ping();err != nil{
		return nil ,err
	}
	return db,nil
}
