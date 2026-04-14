package server

import (
	"ekhoes-server/assets"
	"ekhoes-server/config"
	"ekhoes-server/module"
	"ekhoes-server/utils"
	"html/template"
	"net/http"
	"os"
	"time"
)

var tmpl = template.Must(template.ParseFS(assets.TemplatesFS, "root.htm"))

type RootData struct {
	Package      string
	Version      string
	InstanceName string
	BuildTime    string
	StartTime    string
	UpTime       string
	Container    bool
	Database     string
	Cache        string
	Modules      string
}

func GetRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	/*
	   tmpl, err := template.ParseFiles("template/root.htm")
	   if err != nil {
	       http.Error(w, "File not found", 404)
	       return
	   }
	*/

	//formattedStartTime := config.Runtime.StartTime.UTC().Format("2006-01-02 15:04:05") + " UTC"

	data := RootData{
		Package:      config.Name(),
		Version:      config.Version(),
		InstanceName: os.Getenv("INSTANCE_NAME"),
		BuildTime:    config.BuildTime(),
		StartTime:    config.Runtime.StartTime.Format(time.RFC3339),
		UpTime:       utils.HumanizeDuration(time.Since(config.Runtime.StartTime)),
		Container:    config.IsRunningInContainer(),
		Database:     config.Runtime.Database,
		Cache:        config.Runtime.Cache,
		Modules:      module.GetLoadedModules(),
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}
