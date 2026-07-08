package doctor

import (
	"goxlock/config"
	`time`
)

// - CreateTestUnlockerConfig 
// Testing for unlcoks and its failing points
func (dc *Doctor) CreateTestUnlockerConfig() (error) {
	dc.config.FolderName = dc.config.OutputName
	dc.config.Password = secondarypasswordname
	dc.config.InstructData.Timeout,_ = time.ParseDuration(`1m20s`)
	dc.config.UserAction = config.UnlockFolder
	err := dc.config.Structure()
	if err != nil {
		return err
	}
	return nil
}