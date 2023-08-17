package scripts

import "time"

func ParseStringToTime(s string) (time.Time, error){
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil{
		return time.Now(), err
	}
	return t, nil
}