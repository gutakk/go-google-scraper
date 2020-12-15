package presenters

import "time"

func FormattedCreatedAt(createdAt time.Time) string {
	return createdAt.Format("January 2, 2006")
}
