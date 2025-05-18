package api

import (
	"encoding/json"
	"html/template"
	"net/http"

	"openlocaldog/storage"
)

func HandleGetMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := storage.GetMetrics()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func HandleResetMetrics(w http.ResponseWriter, r *http.Request) {
	storage.ResetMetrics()
	w.WriteHeader(http.StatusNoContent)
}

func HandleGetCardinality(w http.ResponseWriter, r *http.Request) {
	card := storage.GetCardinality()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(card)
}

func HandleDashboard(w http.ResponseWriter, r *http.Request) {
	metrics := storage.GetMetrics()
	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<title>LocalDog Metrics</title>
	<style>
		table { border-collapse: collapse; width: 100%; }
		th, td { padding: 8px; border: 1px solid #ccc; text-align: left; }
	</style>
</head>
<body>
	<h1>LocalDog Metrics</h1>
	<table>
		<tr>
			<th>Name</th><th>Value</th><th>Type</th><th>Tags</th><th>Timestamp</th>
		</tr>
		{{range .}}
		<tr>
			<td>{{.Name}}</td><td>{{.Value}}</td><td>{{.Type}}</td><td>{{range .Tags}}{{.}} {{end}}</td><td>{{.Timestamp}}</td>
		</tr>
		{{end}}
	</table>
</body>
</html>`
	t := template.Must(template.New("dashboard").Parse(tmpl))
	t.Execute(w, metrics)
}
