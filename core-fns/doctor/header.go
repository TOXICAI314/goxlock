package doctor

import "goxlock/config"

// This checks that the integrity program is clearly working or not
func (dc *Doctor) CreateTestVerifyPasswordConfig() (error) {
	dc.config.Password = config.Name
	dc.config.FolderName = dc.config.OutputName
	err := dc.config.Structure()
	if err != nil {
		return err
	}
	return nil 
}