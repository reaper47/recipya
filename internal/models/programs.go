package models

// FFProbe represents the output JSON of the FFprobe program.
type FFProbe struct {
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
}
