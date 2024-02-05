package pixiv

import "fmt"

type PixivRequestError struct {
	StatusCode int
}

func (e *PixivRequestError) Error() string {
	return fmt.Sprintf("upstream returned error %d", e.StatusCode)
}
