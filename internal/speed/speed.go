package speed

type NetSpeed struct {
	Value string
	Unit  string
	Err   error
}

type GetSpeedResp struct {
	DownloadSpeedChannel <-chan NetSpeed
	UploadSpeedChannel   <-chan NetSpeed
}

type SpeedChecker interface {
	GetSpeed() (*GetSpeedResp, error)
}
