package pixiv

type UgoiraApiResponse struct {
	Error   bool
	Message string
	Body    struct {
		Src         string
		OriginalSrc string
		Frames      []struct {
			File  string
			Delay int
		}
	}
}
