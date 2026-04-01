package server
import("bytes";"encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-signalman/internal/store")
func(s *Server)handleList(w http.ResponseWriter,r *http.Request){list,_:=s.db.List();if list==nil{list=[]store.Rule{}};writeJSON(w,200,list)}
func(s *Server)handleCreate(w http.ResponseWriter,r *http.Request){var rule store.Rule;json.NewDecoder(r.Body).Decode(&rule);if rule.Name==""||rule.SourceURL==""||rule.WebhookURL==""{writeError(w,400,"name, source_url, webhook_url required");return};rule.Enabled=true;s.db.Create(&rule);writeJSON(w,201,rule)}
func(s *Server)handleFire(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);var req struct{Payload string `json:"payload"`};json.NewDecoder(r.Body).Decode(&req);rules,_:=s.db.List();var webhookURL string;for _,rl:=range rules{if rl.ID==id{webhookURL=rl.WebhookURL;break}};if webhookURL!=""{go func(){http.Post(webhookURL,"application/json",bytes.NewBufferString(req.Payload))}()};s.db.RecordFire(id,req.Payload);writeJSON(w,200,map[string]string{"status":"fired"})}
func(s *Server)handleToggle(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.Toggle(id);writeJSON(w,200,map[string]string{"status":"toggled"})}
func(s *Server)handleDelete(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.Delete(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
