package doctor

import "goxlock/config"

// This will create a config for the test of the configuration
func (dc *Doctor) CreateTestChangePasswordConfig() (error) {
	dc.config.FolderName = dc.config.OutputName
	dc.config.Password = config.Name
	dc.config.ChangePassword.NewPassword = secondarypasswordname
	dc.config.UserAction = config.ChangePassword
	err := dc.config.Structure()
	if err != nil {
		return err
	}
	return nil
}