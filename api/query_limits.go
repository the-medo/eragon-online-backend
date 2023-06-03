package api

func GetDefaultQueryBoundaries(limit int32, offset int32) (int32, int32) {
	if limit == 0 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	return limit, offset
}
