package main

import (
	"fmt"
	"github.com/civet148/gotools/log"
	"github.com/civet148/sqlca"
	"time"
)

type PhoneCall struct {
	Id                   int64  `db:"id"`
	AccessHash           int64  `db:"access_hash"`
	AdminId              int32  `db:"admin_id"`
	ParticipantId        int32  `db:"participant_id"`
	AdminAuthKeyId       int64  `db:"admin_auth_key_id"`
	ParticipantAuthKeyId int64  `db:"participant_auth_key_id"`
	RandomId             int64  `db:"random_id"`
	AdminProtocol        string `db:"admin_protocol"`
	ParticipantProtocol  string `db:"participant_protocol"`
	GAHash               string `db:"g_a_hash"`
	GA                   string `db:"g_a"`
	GB                   string `db:"g_b"`
	KeyFingerprint       int64  `db:"key_fingerprint"`
	Connections          string `db:"connections"`
	AdminDebugData       string `db:"admin_debug_data"`
	ParticipantDebugData string `db:"participant_debug_data"`
	AdminRating          int32  `db:"admin_rating"`
	AdminComment         string `db:"admin_comment"`
	ParticipantRating    int32  `db:"participant_rating"`
	ParticipantComment   string `db:"participant_comment"`
	Date                 int32  `db:"date"`
	State                int32  `db:"state"`
	CreatedAt            string `db:"created_at"`
	UpdatedAt            string `db:"updated_at"`
}

const (
	TABLE_NAME_PHONE_CALL_SESSIONS = "phone_call_sessions"
)

func main() {

	e := sqlca.NewEngine(true)
	e.Open(sqlca.AdapterSqlx_MySQL, "root:123456@tcp(127.0.0.1:3306)/enterprise?charset=utf8mb4")
	//e.Open(sqlca.AdapterCache_Redis, "redis://127.0.0.1:6379/db?dbnum=0")
	e.Debug(true) //debug on

	var callUpsert = PhoneCall{
		Id:                   0,
		AccessHash:           1234567890,
		AdminId:              1000000,
		ParticipantId:        1000001,
		AdminAuthKeyId:       -6666666431149903665,
		ParticipantAuthKeyId: -7777777424437420153,
		RandomId:             555993992,
		AdminProtocol:        "udp_p2p",
		ParticipantProtocol:  "udp_p2p",
		GAHash:               "746b79e08a1a57868e5e4ed91ebf873c65d668211cf45286048dfcdad0dad8ba",
		GA:                   "",
		GB:                   "",
		KeyFingerprint:       0,
		Connections:          "\"{\"protocol\": \"relay\", \"port\":50001}\"",
		AdminDebugData:       "",
		ParticipantDebugData: "",
		AdminRating:          0,
		AdminComment:         "",
		ParticipantRating:    0,
		ParticipantComment:   "",
		Date:                 0,
		State:                0,
		CreatedAt:            "", //created_at column ignore by insert/upsert/update
		UpdatedAt:            "", //updated_at column ignore by insert/upsert/update
	}

	_ = callUpsert

	var callQuery PhoneCall
	var callList []PhoneCall

	//e.SetPkName("uuid") // set primary key name, default 'id'
	var err error
	var rows int64
	var lastInsertId int64

	_ = lastInsertId

	// insert a record
	lastInsertId, err = e.Model(&callUpsert).Table(TABLE_NAME_PHONE_CALL_SESSIONS).Insert()

	// insert if not exist, otherwise update state and date
	callUpsert.State = 1
	callUpsert.Date = int32(time.Now().Unix())
	lastInsertId, err = e.Model(&callUpsert).Table(TABLE_NAME_PHONE_CALL_SESSIONS).Select("state", "date").Upsert()
	_ = lastInsertId

	//Remark: single record to fetch by primary key which named 'id'
	//SQL: select * from phone_call_sessions where id='1'
	rows, err = e.Model(&callQuery).Table(TABLE_NAME_PHONE_CALL_SESSIONS).
		Id(1).
		Select("id", "access_hash", "admin_id", "participant_id", "admin_auth_key_id", "participant_auth_key_id", "g_a_hash", "created_at", "updated_at").
		Query()
	if err != nil {
		_ = rows
		log.Errorf(err.Error())
		return
	}
	log.Debugf("query result rows [%v] results %+v", rows, log.JsonDebugString(callQuery))

	//Remark: multiple record to fetch by where condition
	//SQL: select id, access_hash, admin_id, participant_id, admin_auth_key_id, participant_auth_key_id from phone_call_sessions where id <='100' limit 5 offset 1
	rows, err = e.Model(&callList).
		Table(TABLE_NAME_PHONE_CALL_SESSIONS).
		//Select("id", "access_hash", "admin_id", "participant_id", "admin_auth_key_id", "participant_auth_key_id", "g_a_hash").
		Where("id <= 100"). // use Where function, the records which be updated can not be refreshed to redis/memcached...
		OrderBy("created_at").
		Desc(). //Asc().
		//GroupBy("admin_id", "participant_id").
		Offset(1).
		Limit(5).
		Query()
	if err != nil {
		_ = rows
		log.Errorf(err.Error())
		return
	}
	log.Debugf("query custom where condition result rows [%v] results %+v", rows, log.JsonDebugString(callList))

	////Remark: single record to fetch by primary key which named 'id', fetch to base type variants
	////SQL: select admin_id, participant_id from phone_call_sessions where id='1'
	var adminId, participantId int64
	rows, err = e.Model(&adminId, &participantId).
		Table(TABLE_NAME_PHONE_CALL_SESSIONS).
		Id(1).
		Select("admin_id", "participant_id").
		Query()
	if err != nil {
		_ = rows
		log.Errorf(err.Error())
		return
	}
	log.Debugf("query result rows [%v] adminId [%d] participantId [%d]", rows, adminId, participantId)

	////Remark: single record to fetch by primary key which named 'id', fetch to base type variants
	////SQL: select id from phone_call_sessions where 1=1 order by id [asc] limit 10 [offset 1]
	var idList []int64
	rows, err = e.Model(&idList).
		Table(TABLE_NAME_PHONE_CALL_SESSIONS).
		Select("id").
		Where("1=1").
		OrderBy("id").
		//Asc().
		Limit(10).
		//Offset(1).
		Query()
	if err != nil {
		_ = rows
		log.Errorf(err.Error())
		return
	}
	log.Debugf("query result rows [%v] id slice %v ", rows, idList)

	//var mapResults = make(map[string]string, 1)
	//rows, err = e.Model(&mapResults).
	//	Table(TABLE_NAME_PHONE_CALL_SESSIONS).
	//	Id(1).
	//	Select("admin_id", "participant_id").
	//	Query()
	//if err != nil {
	//	_ = rows
	//	log.Errorf(err.Error())
	//	return
	//}
	//log.Debugf("query result rows [%v] map [%+v]", rows, mapResults)

	callUpsert.State = 3
	rows, err = e.Model(&callUpsert).
		Table(TABLE_NAME_PHONE_CALL_SESSIONS).
		Id(1).
		Select("state", "connections").
		Update()

	var callRawList []PhoneCall
	strQueryRaw := fmt.Sprintf("SELECT * FROM %v", "phone_call_sessions")
	rows, err = e.Model(&callRawList).QueryRaw(strQueryRaw)
	if err != nil {
		log.Error("QueryRaw error [%v] query [%v]", err.Error(), strQueryRaw)
		return
	}
	log.Debugf("QueryRaw rows [%v] results %+v", rows, callRawList)

	strUpdateRaw := fmt.Sprintf("UPDATE %v SET state='%v' WHERE id='%v'", "phone_call_sessions", 9, 1)
	rows, lastInsertId, err = e.ExecRaw(strUpdateRaw)
	if err != nil {
		log.Error("ExecRaw error [%v] query [%v]", err.Error(), strUpdateRaw)
		return
	}
	log.Debugf("ExecRaw rows [%v] last insert id [%v] query [%v]", rows, lastInsertId, strUpdateRaw)

	//e.Update("123456") //not call Model(), this is a wrong operation, will panic

	log.Info("program exit...")
}
