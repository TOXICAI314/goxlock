package doctor

import (
	"fmt"
	"goxlock/config"
	"goxlock/core-fns/locker"
	"goxlock/core-fns/logger"
	passwordchanger "goxlock/core-fns/password-changer"
	"goxlock/core-fns/profiler"
	"goxlock/core-fns/unlocker"
	"goxlock/core-fns/verify"
	"goxlock/utils"
	"os"
	"time"
)

// This function automatically runs all the comands to check that is there any error formed while the app is running on a device
func (dc *Doctor) Start() error {
	dc.config.SessionID = utils.CreateSessionID() 
	fmt.Printf("Start : Checking Start\n")
	fmt.Println(`--- --- ---`)
	starttime := time.Now()
	err := dc.CreateTestLockerConfig() 
	if err != nil {
		return err
	}
	err = dc.config.Structure()
	if err != nil {
		return err
	}
	err = locker.Locker(&dc.config)
	if err != nil {
		return err
	}
	fmt.Printf("Success : Locking Passed %+v\n",time.Since(starttime))
	fmt.Println(`--- --- ---`)
	starttime = time.Now()
	err = dc.CreateTestVerifyPasswordConfig()
	if err != nil {
		return err
	}
	err = dc.config.Structure()
	if err != nil {
		return err
	}
	err = verify.VerifyUnlock(&dc.config)
	if err != nil {
		return err
	}
	fmt.Printf("Success : Verifying Password Passed %+v\n",time.Since(starttime))
	fmt.Println(`--- --- ---`)
	starttime = time.Now()
	err = dc.CreateTestChangePasswordConfig()
	if err != nil {
		return err
	}
	err = dc.config.Structure()
	if err != nil {
		return err
	}
	err = passwordchanger.ChangePasswordForFolder(&dc.config)
	if err != nil {
		return err
	}
	fmt.Printf("Success : Changing Password Passed %+v\n",time.Since(starttime))
	fmt.Println(`--- --- ---`)
	starttime = time.Now()
	err = dc.CreateTestUnlockerConfig()
	if err != nil {
		return err
	}
	err = dc.config.Structure()
	if err != nil {
		return err
	}
	err = unlocker.Unlocker(&dc.config)
	if err != nil {
		return err
	}
	fmt.Printf("Success : UnLocking Passed %+v\n",time.Since(starttime))
	fmt.Println(`--- --- ---`)
	fmt.Println(`NOTE: You should see a window popup within 2 mins -> Sceduler working`)
	fmt.Println(`--- --- ---`)
	fmt.Println(`All Basic Check parameteres passed!`)
	fmt.Println(`--- --- ---`)
	fmt.Println(`START : Testing for profiler has started`)
	fmt.Println(`--- --- ---`)
	starttime = time.Now()
	pf := &profiler.Profiler{
		Name: `Test`,
		OutputName: `TestfolderOutputName`,
		Instruction: config.Instructions{},
	}
	err = pf.Create()
	if err != nil {
		return err
	}
	fmt.Printf("Success : Profile Creation Passed %+v\n",time.Since(starttime))
	fmt.Println(`--- --- ---`)
	starttime = time.Now()
	err = pf.Update(
		`TestFolderOutputName2`,
		&config.Instructions{},
	)
	if err != nil {
		return err
	}
	fmt.Printf("Success : Profile Updation Passed %+v\n",time.Since(starttime))
	fmt.Println(`--- --- ---`)
	starttime = time.Now()
	err = pf.Delete()
	if err != nil {
		return err
	}
	fmt.Printf("Success : Profile Deletion Passed %+v\n",time.Since(starttime))
	fmt.Println(`--- --- ---`)
	fmt.Println(`All Profile Check parameteres passed!`)
	fmt.Println(`--- --- ---`)
	fmt.Println(`START : Testing for Logger has started`)
	fmt.Println(`--- --- ---`)
	starttime = time.Now()
	logger_x,err := logger.Log(&dc.config,nil)
	if err != nil {
		return err
	}
	err = logger_x.Write()
	if err != nil {
		return err
	}
	fmt.Printf("Success : Logger Creation Passed %+v\n",time.Since(starttime))
	fmt.Println(`--- --- ---`)
	starttime = time.Now()
	formatTime := starttime.Format(`02-01-2006`)
	_,err = logger.ReadLogFile(formatTime)
	if err != nil {
		return err
	}
	fmt.Printf("Success : Logger Read Passed %+v\n",time.Since(starttime))
	err = os.Remove(logger_x.Place)
	if err != nil {
		return &config.FunctionFailError{
			Cause: err.Error(),
			Message: fmt.Sprintf(`Error in removing log file - %s`,logger_x.Place),
			ElapsedTime: time.Now(),
			Provider: `dc.Doctor.Start`,
		}
	}
	fmt.Println(`--- --- ---`)
	fmt.Println(`All Log Check parameteres passed!`)
	fmt.Println(`--- --- ---`)
	
	return nil
}