package model

import (
	"time"
)

type Role_R struct {
	Name struct {
		Th string `json:"th"`
		En string `json:"en"`
	} `json:"name"`
	CreateAtDate   time.Time `json:"createAtDate" bson:"createAtDate"`
	LastupdateDate time.Time `json:"lastupdateDate" bson:"lastupdateDate"`
}

type Char_R struct {
	Name         string    `json:"name"`
	CreateAtDate time.Time `json:"createAtDate" bson:"createAtDate"`
}

type Char_Model struct {
	UID    string `json:"_id" bson:"_id"`
	NAME   string `json:"name" bson:"name"`
	USERID string `json:"userId" bson:"userId"`

	INT     string `json:"INT" bson:"INT"`
	STR     string `json:"STR" bson:"STR"`
	EXP_MIN string `json:"exp_min" bson:"exp_min"`
	EXP_MAX string `json:"exp_max" bson:"exp_max"`
	LEVEL   string `json:"level" bson:"level"`
	IMAGE   string `json:"images" bson:"images"`

	CreateAtDate time.Time `json:"createAtDate" bson:"createAtDate"`
}

type Users_R struct {
	UID              string `json:"_id" bson:"_id"`
	Account_no       string //<< RUNNING NUMBER เชคจาก last index + 1
	Phone            string
	Email            string `json:"email"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	Confirm_password string
	Is_notification  bool
	Is_verify_email  bool
	Active           bool
	Cart             []map[string]interface{} `json:"cart" bson:"cart"`
	Status           struct {
		Online    bool
		LastLogin struct {
			Date      time.Time
			IpAddr    string
			UserAgent string
		}
	}
	Role     string
	RoleID   string
	Block    bool
	Language string `json:"language"`
}

// get user
type Users_R_Agg struct {
	UID        string `json:"_id" bson:"_id"`
	Account_no string `json:"account_no"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Phone      string `json:"phone"`

	Email string `json:"email"`

	Is_notification bool `json:"is_notification"`
	Is_verify_email bool `json:"is_verify_email"`
	Active          bool `json:"active"`
	Status          struct {
		Online    bool `json:"online"`
		LastLogin struct {
			Date      time.Time `json:"date"`
			IpAddr    string    `json:"ip_address"`
			UserAgent string    `json:"user_agent"`
		} `json:"last_login"`
	} `json:"status"`
	Role string `json:"role"`

	Chargingstatus struct {
		ConnectorId   string `json:"connectorId" bson:"connectorId"`
		Firebasetoken string `json:"firebase_token" bson:"firebase_token"`
		ChargerId     string `json:"chargerId" bson:"chargerId"`
	} `json:"charging_status" bson:"charging_status"`

	StartWorkingAtDate time.Time `json:"start_working_at_date"`
	CreateAtDate       time.Time `json:"create_at_date"`
}
