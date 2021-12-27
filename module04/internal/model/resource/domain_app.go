package resource

import (
	"github.com/pkg/errors"
	"self/internal/dao"
	"time"
)

type DomainApp struct {
	Id      int64     `json:"id"`
	Domain  string    `json:"domain"`
	AppId   string    `json:"app_id"`
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

func (this *DomainApp) GetsAll() ([]*DomainApp, error) {
	ret := make([]*DomainApp, 0)
	err := dao.Db().Find(&ret)
	return ret, err
}

func (this *DomainApp) RemoveAndInsert(data []*DomainApp) error {
	session := dao.Db().NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		if eR := session.Rollback(); eR != nil {
			return errors.New(eR.Error())
		}
		return errors.New(err.Error())
	}
	_, err := session.Where("id > ?", 0).Delete(new(DomainApp))
	if err != nil {
		if eR := session.Rollback(); eR != nil {
			return errors.New(eR.Error())
		}
		return errors.New(err.Error())
	}
	var length = 500
	for i := 0; i < len(data); i += length {
		if i+length > len(data) {
			if _, err = session.Insert(data[i:]); err != nil {
				if eR := session.Rollback(); eR != nil {
					return errors.New(eR.Error())
				}
				return errors.New(err.Error())
			}
		} else {
			if _, err = session.Insert(data[i : i+length]); err != nil {
				if eR := session.Rollback(); eR != nil {
					return errors.New(eR.Error())
				}
				return errors.New(err.Error())
			}
		}
	}
	if err = session.Commit(); err != nil {
		return errors.New(err.Error())
	}
	return nil
}
