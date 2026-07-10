package cmd

import (
	"encoding/json"
	"fmt"
	"goxlock/config"
	"goxlock/core-fns/doctor"
	"goxlock/core-fns/logger"
	"goxlock/core-fns/performer"
	"goxlock/core-fns/profiler"
	"goxlock/core-fns/updater"
	"goxlock/utils"
	"time"

	"github.com/spf13/cobra"
)

// - Helpful Variables
// These varibales will only be used in cmd/root.go for assists
var (
	provider 	string  	= 		`cmd.root`
)

// - Binded Variables
// Info : Binded variable are binded towards a flag which dumps its value to the variable
var (
	// FolderName -> The main variable that contains the folder info of the user command
	FolderName string
	// OutputName -> The name by which the user want the output to be
	OutputName string
	// Password -> The secrent token by which the folder will be encrypted
	Password string
)

// - Hidden Bind Variables
// Used for the application internal use and lower case uses
var (
	// relock -> switch to relock the file after a given duration
	relock bool
	// sessionId -> Gets the session by which relock file is made and will give the necessary data
	sessionId string
)

// - Trivial Acting
// - Variables that are just needed to activate certain features and plays no role afterwards
var (
	// - Config Section
	// Info : These variabeles are to setup actions for the config

	// lock -> Triggers lock action
	lock 			bool
	// unlock -> Triggers unlock action
	unlock 			bool
	// changePassword -> Triggers password changing action
	changePassword 	bool
	// verifyPassword -> Triggers verify password action
	verifyPassword 	bool
	// header -> Triggers header seeking mechanism
	header 			bool

	// - Instruction Section
	// Info : These variables are for the instructions that comes with config

	// deloriginal -> Triggers self deletion of the folder or the file
	deleteoriginal bool
	// timeoutS -> Triggers the timeout mechanism which then calls the scheduler
	timeoutS string
	// exclude -> Excludes the file pattern that user dont want
	exlude []string
	// stats -> Toggles the stats of the action as needed
	stats bool
	// unsafe -> Removes the safety bearing of the application
	unsafe bool

	// loggerallowed -> Toggles the logger that will write the logs into the log file
	loggerallowed bool

	// - Profile
	// Info : Deals with the arguments passed for the profiling system

	// profileName -> profiles the name of the current profile
	profileName string

	// - Log
	// Deals with the argumenst passe for the loging section
	logdate 	string

	// - Updater
	// Deals with update and stuff
	checkupdate bool

	// - TUI
	// The triggering point of tui
	serve_tui 		bool
)

// - One time runners
// These will run at the starting of the code if given and will return immediately without causeing any action
var (
	// doctorcheck -> fires the checkup for the perfect working of the application on your device
	doctorcheck 		bool
	// readlog -> tells to read the log from the given timed log file
	readlog 			bool
	// makeprofile -> fires the running of the profile making mechanism
	makeprofile 		bool
	// deleteprofile -> deletes the profile that is named
	deleteprofile	 	bool
	// updateprofile -> updates the feilds of the profile user wanted to use
	updateprofile 		bool
	// useprofile -> Intead of returning -> It let the command run further
	useprofile 			bool
	// where -> Just gives the current working place to the caller
	where 				bool
)

// GlobalCfg : For all the user request passed into sigularity
var Cfg *config.Config = &config.Config{}

// - COMMAND DEFINE

// rootcmd -> is the main command that is used to trigger all the tasks
var rootcmd *cobra.Command = &cobra.Command{
	// - Base definition
	// These will be used to run the main command and get the version info
	// Any change in `Use` field -> Different name for running the command
	Use:     config.Name,
	Long: 	 config.Banner,
	Version: utils.GetBasicVersionDetails(fmt.Sprintf(`%.1f`,config.Version)),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// - PreRunner
		// This will run in the starting doing the objective that is assigned to it or left by previous work

		// Info : This will clear the cache that is stored into the temporary folder
		err := utils.ClearAllTempFolderJunk(config.Name)
		if err != nil {
			fmt.Println(&config.FunctionFailError{
				Cause:   err.Error(),
				Message: fmt.Sprintf(`Cannot make clear up the previous temp folder junk made by %s`, config.Name),
				ElapsedTime: time.Now(),
				Provider: `utils.ClearAllTempFolderJunk`,
			})
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {

		// - One time Fire commands /
		// These will fire and return without encroaching the action
		switch {
		case checkupdate:
			return updater.CheckForUpdate()
		case where:
			path,err := utils.Where()
			if err != nil {
				return &config.FunctionFailError{
					Cause: err.Error(),
					Message: `The underlying function cant find the executable path of the current running program`,
					ElapsedTime: time.Now(),
					Provider: `utils.Where`,
				}
			}
			fmt.Printf(`Current Woking path : %s`,path)
			return nil
		// This will check for any doctor check that is needed (NOTE : Doctor only checks for functionality that the application wants not the functions)
		case doctorcheck:
			dc := &doctor.Doctor{}
			return dc.Start() 
		case readlog:
			logger, err := logger.ReadLogFile(logdate)
			if err != nil {
				return err
			}
			for _, data := range logger.CollectiveLoggerData.Logs {
				data, err := json.MarshalIndent(data, ``, `	`)
				if err != nil {
					continue
				}
				fmt.Println((string(data)))
				fmt.Println(`--- --- ---`)
			}
			return nil
		// Makes the profile from the user request by utilising the already provided flags
		case makeprofile:
			var err error
			var timeout time.Duration
			if timeoutS != `` {
				timeout, err = time.ParseDuration(timeoutS)
				if err != nil {
					return &config.FunctionCancelError{
						Cause:   err.Error(),
						Message: fmt.Sprintf(`Cannot parse the given time : %s`,timeoutS),
						ElapsedTime: time.Now(),
						Provider: `time.ParseDuration`,
					}
				}
			}
			if OutputName == `` {
				return &config.FunctionCancelError{
					Cause:   `Invalid output name`,
					Message: `Provide a Valid name to be taken for the profile`,
					Provider: provider,
					ElapsedTime: time.Now(),
				}
			}
			pf := &profiler.Profiler{
				Name:       profileName,
				OutputName: OutputName,
				Instruction: config.Instructions{
					Timeout:        timeout,
					DeleteOriginal: deleteoriginal,
					Exclusion:      exlude,
					Stats:          stats,
					LoggerAllowed:  loggerallowed,
					UnSafe:         unsafe,
				},
			}
			err = pf.Create()
			if err == nil {
				fmt.Printf("Creation of the Profile - %s Successful", pf.Name)
			}
			return err

		// Deletes those profiles which are named by the user
		// It deletes only one profile at a time
		case deleteprofile:
			pf := &profiler.Profiler{
				Name: profileName,
			}
			err := pf.Delete()
			if err == nil {
				fmt.Printf("Deletion of the Profile - %s Successful", pf.Name)
			}
			return err

		// Updates the profile by getting all the inputs -> Fetching all the data -> Rewritting a new set of data
		case updateprofile:
			pf := &profiler.Profiler{
				Name: profileName,
			}
			instr := &config.Instructions{
				DeleteOriginal: deleteoriginal,
				Exclusion:      exlude,
				Stats:          stats,
				LoggerAllowed:  loggerallowed,
				UnSafe:         unsafe,
			}
			var err error
			if timeoutS != `` {
				instr.Timeout, err = time.ParseDuration(timeoutS)
				if err != nil {
					return &config.FunctionFailError{
						Cause:   err.Error(),
						Message: `Bad time input which cant be parsed`,
						Provider: `time.ParseDuration`,
						ElapsedTime: time.Now(),
					}
				}
			}
			err = pf.Update(OutputName, instr)
			if err == nil {
				fmt.Printf("Updation of the Profile - %s Successful", pf.Name)
			}
			return err

		// Instead of returning it uses the profile and names the variables that are going to be used by its returned data
		case useprofile:
			pf := &profiler.Profiler{
				Name: profileName,
			}
			err := pf.Fetch()
			if err != nil {
				return err
			}
			timeoutS = pf.Instruction.Timeout.String()
			OutputName = pf.OutputName
			deleteoriginal = pf.Instruction.DeleteOriginal
			timeoutS = pf.Instruction.Timeout.String()
			exlude = pf.Instruction.Exclusion
			unsafe = pf.Instruction.UnSafe
		}
		// - Incompatible
		// Incomapttible commands which may conflict in their works will get struck out here
		switch {
		case verifyPassword, changePassword, lock, unlock, header:
			if FolderName == `` {
				return &config.FunctionCancelError{
					Cause:   `Together operations cant be neglected -> 'lock|unlock|changePassword|verifyPassword|header' && 'folder'`,
					Message: `Provide Sufficient details for verifying`,
					ElapsedTime: time.Now(),
					Provider: provider,
				}
			}
		case relock:
			if sessionId == `` {
				return &config.FunctionCancelError{
					Cause:   `Together operations are neglected -> 'relock' && 'sessionID'`,
					Message: `Together coming operations relocking`,
					ElapsedTime: time.Now(),
					Provider: provider,
				}
			}
		}

		// - Data Dump
		// - Info : All the user request will be dumped into the global config struct `Cfg`

		// - Config Section

		var Timeout time.Duration
		var err error
		var ans bool
		if timeoutS != `` {
			if !unsafe {
				// - Confirmation
				// Info : As info stores the os level enctypted password directly to the disk
				// If any attacker gets in as the user , they can decrypt the file while the `session.json` file is there
				ans = utils.GetYesORNo(
					fmt.Sprintf(`
				SECURITY ALLERT:
				
				- By using 'time-out' you are allowing the app to store the encrypted data into the disk , which may contain your password.
				- Password is not in raw string or bytes but is encrypted.
				- But for the duration of %s that file will be stored into your disk and is vulnarable from any attacker if they get the os level id as you.


				Better lock your device so that no one from your side can decrypt your data for that time period.
				As once the scheduler has done working , the scheduled file will be auto deleted.

				Are you sure you wanna schedule this task OR Just do the action without any timeout ? (y/n) : `, timeoutS),
				)
			}
			if !ans {
				Timeout, err = time.ParseDuration(timeoutS)
				if err != nil {
					return &config.FunctionFailError{
						Cause:   err.Error(),
						Message: `Bad time input which cant be parsed`,
						Provider: `time.ParseDuration`,
						ElapsedTime: time.Now(),
					}
				}
			}	
		}
		var action int

		// - Singularity Check
		// Info : This checks for the passing of the only action , as two or more can result in undefined behaviour

		var actionArray []bool = make([]bool, 0)
		actioncount := 0
		actionArray = append(actionArray, lock, unlock, relock, changePassword, verifyPassword, header)
		for _, a := range actionArray {
			if a {
				actioncount++
			}
		}
		if actioncount > 1 {
			return &config.FunctionCancelError{
				Cause: `More action than one`,
				Message: `Your current action tried to do more than one action
				Provided actions :
				- Lock
				- Unlock
				- Relock
				- Change Password
				- Verify Password
				- Header
				`,
				Provider: provider,
				ElapsedTime: time.Now(),
			}
		}

		switch {
		case lock:
			action = config.LockFolder
		case unlock:
			action = config.UnlockFolder
		case relock:
			action = config.RelockFolder
		case changePassword:
			action = config.ChangePassword
		case verifyPassword:
			action = config.VerifyPassword
		case header:
			action = config.HeaderCheck
		}
		Cfg.UserAction = action
		Cfg.FolderName = FolderName
		Cfg.SessionID = sessionId
		Cfg.OutputName = OutputName
		Cfg.StartedAt = time.Now()

		// - Password Security
		// Info : Gets the password from the user and hide it from the terminal buffer
		// As terminal bugger stores every stuff , its neccessary to delete whats need to be private

		switch Cfg.UserAction {
		case config.LockFolder:
			Password, err = utils.ConfirmSecretPassword(`Enter your password: `)
			if err != nil {
				return &config.FunctionFailError{
					Cause:   err.Error(),
					Message: `There is an internal error while recording the password`,
					ElapsedTime: time.Now(),
					Provider: provider,
				}
			}
			if Password == `` {
				return &config.FunctionCancelError{
					Cause:   `Empty Password`,
					Message: `Cannot Continue with an empty password`,
					ElapsedTime: time.Now(),
					Provider: provider,
				}
			}
		case config.UnlockFolder:
			Password, err = utils.SecretInput(`Enter your password: `)
			if err != nil {
				return &config.FunctionFailError{
					Cause:   err.Error(),
					Message: `Cannot create a secure Path to read the password ; Make sure you have opened the app in secure terminal`,
					ElapsedTime: time.Now(),
					Provider: provider,
				}
			}
		case config.ChangePassword:
			Password, err = utils.SecretInput(`Enter your initial passwrord: `)
			if err != nil {
				return &config.FunctionFailError{
					Cause:   err.Error(),
					Message: `Cannot create a secure Path to read the password ; Make sure you have opened the app in secure terminal`,
					ElapsedTime: time.Now(),
					Provider: provider,
				}
			}

			newPassword, err := utils.ConfirmSecretPassword(`Enter your new password: `)
			if err != nil {
				return &config.FunctionFailError{
					Cause:   err.Error(),
					Message: `There is an internal error while recording the password`,
					ElapsedTime: time.Now(),
					Provider: provider,
				}
			}
			if newPassword == `` {
				return &config.FunctionCancelError{
					Cause:   `Empty Password`,
					Message: `Cannot Continue with an empty password`,
					ElapsedTime: time.Now(),
					Provider: provider,
				}
			}
			Cfg.ChangePassword = config.ChangePasswordData{
				NewPassword: newPassword,
			}
		case config.VerifyPassword:
			Password, err = utils.SecretInput(`Enter your password: `)
			if err != nil {
				return &config.FunctionFailError{
					Cause:   err.Error(),
					Message: `Cannot create a secure Path to read the password ; Make sure you have opened the app in secure terminal`,
					ElapsedTime: time.Now(),
					Provider: provider,
				}
			}
		}
		Cfg.Password = Password

		// - Instruction Section
		Cfg.InstructData.DeleteOriginal = deleteoriginal
		Cfg.InstructData.Timeout = Timeout
		Cfg.InstructData.Exclusion = exlude
		Cfg.InstructData.Stats = stats
		Cfg.InstructData.LoggerAllowed = loggerallowed
		Cfg.InstructData.UnSafe = unsafe

		err = Cfg.Structure()
		if err != nil {
			return err
		}

		// - Actions
		return performer.PerformAction(Cfg)
	},
}

/*
Execute : The function that will tie itself to the `main` function
- Allies of the `rootcmd.Execute()`
*/
func Execute() error {
	return rootcmd.Execute()
}

// - Init function
// 1. init() -> Runs in the startup of this file and will collect all the values binded to the flags
func init() {

	// - Usable Flags
	// The Flags that are not hidden and is visible to all the users
	rootcmd.Flags().StringVarP(&FolderName, `folder`, `f`, ``, `Folder to encrypt`)
	rootcmd.Flags().BoolVar(&lock, `lock`, false, `Lock the given folder`)
	rootcmd.Flags().BoolVar(&deleteoriginal, `del-original`, false, `Deletes the original folder`)
	rootcmd.Flags().StringVar(&OutputName, `out`, ``, `Gives the custom name to the output folder`)
	rootcmd.Flags().BoolVar(&unlock, `unlock`, false, `Unlocks the targetted folder`)
	rootcmd.Flags().StringVar(&timeoutS, `time-out`, ``, `Shedule re encryption over a certain time (must give minute/hour/day timelimits )`)
	rootcmd.Flags().BoolVar(&changePassword, `change-password`, false, `Changes your current password`)
	rootcmd.Flags().BoolVar(&verifyPassword, `verify-password`, false, `Verify your folder password without unlocking it`)
	rootcmd.Flags().StringSliceVar(&exlude, `exclude`, []string{}, `Excludes given patterns`)
	rootcmd.Flags().BoolVar(&loggerallowed, `log`, false, `Allow logger to log into the hardcoded files`)
	rootcmd.Flags().BoolVar(&header, `header`, false, `Check header data of the encrypted file`)
	rootcmd.Flags().BoolVar(&stats, `stats`, false, `Use it to see the stats of your operation`)
	rootcmd.Flags().BoolVar(&readlog, `read-log`, false, `Toggle the read log function`)
	rootcmd.Flags().StringVar(&logdate, `log-date`, ``, `Gives the log date to instructor`)
	rootcmd.Flags().BoolVar(&where,`where`,false,`Gives current working executable path`)

	// - Profile
	rootcmd.Flags().BoolVar(&makeprofile, `make-profile`, false, `Makes the profile of the user as per name`)
	rootcmd.Flags().StringVar(&profileName, `profile-name`, ``, `Stores the name of the profile`)
	rootcmd.Flags().BoolVar(&deleteprofile, `del-profile`, false, `Deletes the named profile`)
	rootcmd.Flags().BoolVar(&updateprofile, `update-profile`, false, `Updates the named profile with the given data`)
	rootcmd.Flags().BoolVar(&useprofile, `use-profile`, false, `Use the saved profile`)

	// - Updater
	rootcmd.Flags().BoolVar(&checkupdate, `check-update`, false, `Checks for provided updates from trusted sources`)

	// - Hidden Flags
	// Hidden Flags -> For super users (features of future) or application internal use calls
	rootcmd.Flags().StringVar(&sessionId, `session`, ``, `Use it to assign a session to the folder`)
	rootcmd.Flags().BoolVar(&relock, `re-lock`, false, `Toggle relocking mechanism`)
	rootcmd.Flags().BoolVar(&doctorcheck, `doctor`, false, `Check the internal parts of the code`)
	rootcmd.Flags().BoolVar(&unsafe, `unsafe`, false, `Makes the current flow unsafe`)

	// Actual Hiding ->
	rootcmd.Flags().MarkHidden(`re-lock`)
	rootcmd.Flags().MarkHidden(`session`)
	rootcmd.Flags().MarkHidden(`doctor`)
}
