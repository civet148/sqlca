package sqlca

import "testing"

const (
	TableNameUsers = "users"
)

const (
	USER_COLUMN_ID           = "id"
	USER_COLUMN_NAME         = "name"
	USER_COLUMN_PHONE        = "phone"
	USER_COLUMN_SEX          = "sex"
	USER_COLUMN_EMAIL        = "email"
	USER_COLUMN_DISABLE      = "disable"
	USER_COLUMN_BALANCE      = "balance"
	USER_COLUMN_SEX_NAME     = "sex_name"
	USER_COLUMN_DATA_SIZE    = "data_size"
	USER_COLUMN_EXTRA_DATA   = "extra_data"
	USER_COLUMN_CREATED_TIME = "created_time"
	USER_COLUMN_UPDATED_TIME = "updated_time"
)

type UserDO struct {
	Id          uint64  `json:"id" db:"id" bson:"_id"`                                               //auto inc id
	Name        string  `json:"name" db:"name" bson:"name"`                                          //user name
	Phone       string  `json:"phone" db:"phone" bson:"phone"`                                       //phone number
	Sex         uint8   `json:"sex" db:"sex" bson:"sex"`                                             //user sex
	Email       string  `json:"email" db:"email" bson:"email"`                                       //email
	Disable     int8    `json:"disable" db:"disable" bson:"disable"`                                 //disabled(0=false 1=true)
	Balance     Decimal `json:"balance" db:"balance" bson:"balance"`                                 //balance of decimal
	SexName     string  `json:"sex_name" db:"sex_name" bson:"sex_name"`                              //sex name
	DataSize    int64   `json:"data_size" db:"data_size" bson:"data_size"`                           //data size
	ExtraData   string  `json:"extra_data" db:"extra_data" sqlca:"isnull" bson:"extra_data"`         //extra data (json)
	CreatedTime string  `json:"created_time" db:"created_time" sqlca:"readonly" bson:"created_time"` //created time
	UpdatedTime string  `json:"updated_time" db:"updated_time" sqlca:"readonly" bson:"updated_time"` //updated time
}

var db *Engine

func init() {
	var err error
	db, err = NewEngine("mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4")
	if err != nil {
		panic(err.Error())
	}
}

func TestPrepareStatement(t *testing.T) {

}
