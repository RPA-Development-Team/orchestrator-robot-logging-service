package entity

type Log struct {
	ID           uint64 `json:"id" gorm:"primary_key;auto_increment"`
	LogType      string `json:"logType"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	Timestamp    string `json:"timestamp"`
	Message      string `json:"message"`
	RobotAddress string `json:"robotAddress"`
	UserId       string `json:"userId"`
}
