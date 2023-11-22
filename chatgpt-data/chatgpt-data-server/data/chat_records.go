package data

import (
	"chatgpt-data/pkg/log"
	"database/sql"
	"strings"
)

type IChatRecordsData interface {
	Add(record *ChatRecord) error
	GetById(id int64) (record *ChatRecord, err error)
}

type ChatRecord struct {
	ID              int64    `json:"id"`
	Account         string   `json:"account"`
	GroupID         string   `json:"group_id"`
	UserMsg         string   `json:"user_msg"`
	UserMsgTokens   int      `json:"user_msg_tokens"`
	UserMsgKeywords []string `json:"user_msg_keywords"`
	AIMsg           string   `json:"ai_msg"`
	AIMsgTokens     int      `json:"ai_msg_tokens"`
	ReqTokens       int      `json:"req_tokens"`
	CreateAt        int64    `json:"create_at"`
	EnterpriseId    string   `json:"enterprise_id"`
	Endpoint        int      `json:"endpoint"`
	EndpointAccount string   `json:"endpoint_account"`
}

type chatRecordsData struct {
	log log.ILogger
	db  *sql.DB
}

func NewChatRecordsData(db *sql.DB, log log.ILogger) IChatRecordsData {
	return &chatRecordsData{
		log: log,
		db:  db,
	}
}

func (data *chatRecordsData) Add(cr *ChatRecord) (err error) {
	sqlStr := "insert into chat_records(account,group_id,user_msg,user_msg_tokens,user_msg_keywords,ai_msg,ai_msg_tokens,req_tokens,create_at,enterprise_id,endpoint,endpoint_account)values(?,?,?,?,?,?,?,?,?,?,?,?)"
	res, err := data.db.Exec(sqlStr, cr.Account, cr.GroupID, cr.UserMsg, cr.UserMsgTokens, strings.Join(cr.UserMsgKeywords, ","), cr.AIMsg, cr.AIMsgTokens, cr.ReqTokens, cr.CreateAt, cr.EnterpriseId, cr.Endpoint, cr.EndpointAccount)
	if err != nil {
		data.log.Error(err)
		return
	}
	cr.ID, _ = res.LastInsertId()
	return
}
func (data *chatRecordsData) GetById(id int64) (cr *ChatRecord, err error) {
	sqlStr := "select id,account,group_id,user_msg,user_msg_tokens,user_msg_keywords,ai_msg,ai_msg_tokens,req_tokens,create_at,enterprise_id,endpoint,endpoint_account from chat_records where id = ?"
	row := data.db.QueryRow(sqlStr, id)
	cr = &ChatRecord{}
	var keywords string
	err = row.Scan(&cr.ID, &cr.Account, &cr.GroupID, &cr.UserMsg, &cr.UserMsgTokens, &keywords, &cr.AIMsg, &cr.AIMsgTokens, &cr.ReqTokens, &cr.CreateAt, &cr.EnterpriseId, &cr.Endpoint, &cr.EndpointAccount)
	if err != nil {
		data.log.Error(err)
		return nil, err
	}
	cr.UserMsgKeywords = strings.Split(keywords, ",")
	return cr, err
}
