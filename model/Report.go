package model

type Report struct {
	Id     int    `json:"id" gorm:"primaryKey"`
	BlogId int    `json:"blogId"`
	UserId int    `json:"userId"`
	Reason string `json:"reason"`
}

func NewReport(userId int, blogId int, reason string) (*Report, error) {
	report := &Report{
		UserId: userId,
		BlogId: blogId,
		Reason: reason,
	}
	return report, nil
}
