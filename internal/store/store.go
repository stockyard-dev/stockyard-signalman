package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type AlertRule struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Condition string `json:"condition"`
	Threshold int `json:"threshold"`
	Channel string `json:"channel"`
	Target string `json:"target"`
	Enabled int `json:"enabled"`
	FireCount int `json:"fire_count"`
	LastFiredAt string `json:"last_fired_at"`
	Cooldown int `json:"cooldown"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"signalman.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS alert_rules(id TEXT PRIMARY KEY,name TEXT NOT NULL,condition TEXT DEFAULT '',threshold INTEGER DEFAULT 0,channel TEXT DEFAULT 'webhook',target TEXT DEFAULT '',enabled INTEGER DEFAULT 1,fire_count INTEGER DEFAULT 0,last_fired_at TEXT DEFAULT '',cooldown INTEGER DEFAULT 300,created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *AlertRule)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO alert_rules(id,name,condition,threshold,channel,target,enabled,fire_count,last_fired_at,cooldown,created_at)VALUES(?,?,?,?,?,?,?,?,?,?,?)`,e.ID,e.Name,e.Condition,e.Threshold,e.Channel,e.Target,e.Enabled,e.FireCount,e.LastFiredAt,e.Cooldown,e.CreatedAt);return err}
func(d *DB)Get(id string)*AlertRule{var e AlertRule;if d.db.QueryRow(`SELECT id,name,condition,threshold,channel,target,enabled,fire_count,last_fired_at,cooldown,created_at FROM alert_rules WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.Condition,&e.Threshold,&e.Channel,&e.Target,&e.Enabled,&e.FireCount,&e.LastFiredAt,&e.Cooldown,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]AlertRule{rows,_:=d.db.Query(`SELECT id,name,condition,threshold,channel,target,enabled,fire_count,last_fired_at,cooldown,created_at FROM alert_rules ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []AlertRule;for rows.Next(){var e AlertRule;rows.Scan(&e.ID,&e.Name,&e.Condition,&e.Threshold,&e.Channel,&e.Target,&e.Enabled,&e.FireCount,&e.LastFiredAt,&e.Cooldown,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM alert_rules WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM alert_rules`).Scan(&n);return n}
