package entity

type Log struct {
	ID        uint64 `json:"id" gorm:"primary_key;auto_increment"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
	RobotID   uint64 `json:"robotId"`
}
