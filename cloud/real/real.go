package real

import "time"

type Tweet struct {
	ID         int64
	Text       string
	Favourites int
	Retweets   int
	Date       time.Time
	Device     string

	IsRetweet bool
	IsDeleted bool

	IsFlagged bool

	//not sure if we want this for now
	Real bool
}

func RandomSample() {

}
