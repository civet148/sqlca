package main

import (
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
}

const (
	TABLE_NAME_PHONE_CALL_SESSIONS = "phone_call_sessions"
)

func main() {

	e := sqlca.NewEngine(true)
	e.Open(sqlca.AdapterSqlx_MySQL, "root:123456@tcp(127.0.0.1:3306)/enterprise?charset=utf8mb4")
	//e.Open(sqlca.AdapterCache_Redis, "redis://127.0.0.1:6379/db?dbnum=0")

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
		Connections:          "",
		AdminDebugData:       "",
		ParticipantDebugData: "",
		AdminRating:          0,
		AdminComment:         "",
		ParticipantRating:    0,
		ParticipantComment:   "",
		Date:                 0,
		State:                0,
	}

	var callQuery PhoneCall
	var callList []PhoneCall

	// insert a record
	id, err := e.Model(&callUpsert).Table(TABLE_NAME_PHONE_CALL_SESSIONS).Insert()
	_ = id

	// insert if not exist, otherwise update state and date
	callUpsert.State = 1
	callUpsert.Date = int32(time.Now().Unix())
	id, err = e.Model(&callUpsert).Table(TABLE_NAME_PHONE_CALL_SESSIONS).Upsert("state", "date")
	_ = id

	//Remark: single record to fetch by primary key which named 'id'
	//SQL: select * from phone_call_sessions where id='99'
	var rows int64
	rows, err = e.Model(&callQuery).Table(TABLE_NAME_PHONE_CALL_SESSIONS).
		Id(1).
		//Select("id", "access_hash", "admin_id", "participant_id", "admin_auth_key_id", "participant_auth_key_id", "g_a_hash").
		Query()
	if err != nil {
		_ = rows
		log.Errorf(err.Error())
		return
	}
	log.Debugf("query result rows [%v] values %+v", rows, callQuery)

	//Remark: multiple record to fetch by where condition
	//SQL: select id, access_hash, admin_id, participant_id, admin_auth_key_id, participant_auth_key_id from phone_call_sessions where id <='100'
	rows, err = e.Model(&callList).
		Table(TABLE_NAME_PHONE_CALL_SESSIONS).
		//Select("id", "access_hash", "admin_id", "participant_id", "admin_auth_key_id", "participant_auth_key_id", "g_a_hash").
		Where("id <= 100"). // use Where function, the records which be updated can not be refreshed to redis/memcached...
		Query()
	if err != nil {
		_ = rows
		log.Errorf(err.Error())
		return
	}
	log.Debugf("query result rows [%v] values %+v custom where condition", rows, callList)

	callUpsert.State = 3
	rows, err = e.Model(&callUpsert).Table(TABLE_NAME_PHONE_CALL_SESSIONS).Id(1).Update("state")
	log.Debugf("update effect rows [%v] ", rows)

	log.Info("program exit...")
}
