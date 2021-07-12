package storage

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/venus-auth/config"
	"github.com/filecoin-project/venus-auth/core"
	"github.com/filecoin-project/venus-auth/log"
	"strings"
	"time"
)

func NewStore(cnf *config.DBConfig, dataPath string) (Store, error) {
	switch strings.ToLower(cnf.Type) {
	case config.Mysql:
		log.Warn("mysql storage")
		return newMySQLStore(cnf)
	case config.Badger:
		log.Warn("badger storage")
		return newBadgerStore(dataPath)
	}
	return nil, fmt.Errorf("the type %s is not currently supported", cnf.Type)
}

type Store interface {
	Put(kp *KeyPair) error
	Delete(token Token) error
	Has(token Token) (bool, error)
	List(skip, limit int64) ([]*KeyPair, error)

	// user
	HasUser(name string) (bool, error)
	GetUser(name string) (*User, error)
	HasMiner(maddr address.Address) (bool, error)
	GetMiner(maddr address.Address) (*User, error)
	PutUser(*User) error
	UpdateUser(*User) error
	ListUsers(skip, limit int64, state int, sourceType core.SourceType, code core.KeyCode) ([]*User, error)
}

type KeyPair struct {
	Name       string    `gorm:"column:name;type:varchar(50);NOT NULL"`
	Perm       string    `gorm:"column:perm;type:varchar(50);NOT NULL"`
	Extra      string    `gorm:"column:extra;type:varchar(255);"`
	Token      Token     `gorm:"column:token;type:varchar(512);uniqueIndex:token_token_IDX,type:hash;not null"`
	CreateTime time.Time `gorm:"column:createTime;type:datetime;NOT NULL"`
}

func (*KeyPair) TableName() string {
	return "token"
}

type Token string

func (t Token) Bytes() []byte {
	return []byte(t)
}
func (t Token) String() string {
	return string(t)
}

func (t *KeyPair) CreateTimeBytes() ([]byte, error) {
	val, err := t.CreateTime.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (t *KeyPair) FromBytes(key []byte, val []byte) error {
	t.Token = Token(key)
	tm := time.Time{}
	err := tm.UnmarshalBinary(val)
	if err != nil {
		return err
	}
	t.CreateTime = tm
	return nil
}

type User struct {
	Id         string          `gorm:"column:id;type:varchar(255);primary_key"`
	Name       string          `gorm:"column:name;type:varchar(50);uniqueIndex:users_name_IDX,type:btree;not null"`
	Miner      string          `gorm:"column:miner;type:varchar(255);index:users_miner_IDX,type:btree"`
	Comment    string          `gorm:"column:comment;type:varchar(255);"`
	SourceType core.SourceType `gorm:"column:stype;type:tinyint(4);default:0;NOT NULL"`
	State      int             `gorm:"column:state;type:tinyint(4);default:0;NOT NULL"`
	ReqLimit   ReqLimit        `gorm:"column:reqLimit;type:varchar(512)"`
	CreateTime time.Time       `gorm:"column:createTime;type:datetime;NOT NULL"`
	UpdateTime time.Time       `gorm:"column:updateTime;type:datetime;NOT NULL"`
}

type ReqLimit struct {
	Cap      int64
	ResetDur time.Duration
}

func (rl *ReqLimit) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	if len(bytes) == 0 {
		*rl = ReqLimit{}
		return nil
	}
	return json.Unmarshal(bytes, rl)
}

func (rl ReqLimit) Value() (driver.Value, error) {
	return json.Marshal(rl)
}

func (*User) TableName() string {
	return "users"
}

func (t *User) Bytes() ([]byte, error) {
	buff, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return buff, nil
}

func (t *User) FromBytes(buff []byte) error {
	err := json.Unmarshal(buff, t)
	return err
}

func (t *User) CreateTimeBytes() ([]byte, error) {
	val, err := t.CreateTime.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return val, nil
}
