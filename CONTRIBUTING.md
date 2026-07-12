# WELCOME

Thank you for your interest in contributing to `goxlock`.
Whether it's a bug report, documentation improvement, or a new feature, every contribution is appreciated.

# Info

## Version

`v1.0.0.0-Alpha`

## Email

`imranjmi008@gmail.com`

## Quick Help Lookup

If you can see that you can do any from this list you can read rules go to the bottom section . Issues are discussed there.

- Logic golang perogramming
- Documentation Improvement.
- Performance Optimizations
- TUI
- Linux Support

# PRECAUTIONS

Given are the set of rules that you must read before contributing to the program. These will lay down what to do and how to do so that the issue is solved efficiently.

## Before Opening an Issue

There must be certain precaution that needs to be taken before opening an issue for the `goxlock`:

- Search existing issues and confirm the problem already exists or not.
- Confirm the version that the issue was reported to.

- If you see an issue you can tell to the provided emails or raise it , it will be shown or be fixed in the next session.

## Coding Guidlines

These followig points will lay down the structure by which your contribution must abide:

- Your function must go from low level -> high level (A base function may not handle everything).
- Form Different files for different sections of function.
- Panics are not allowed in the runtime
- Error must be returned in form of `FunctionCancelError` -> When function cuts itself because of a constraint not of error  &  `FunctionFailError` -> When function ends because of an error .
- Apply these two error struct in base function call not the parent caller , other wise it will create an ugly nested structure.
- Must not disturb other Action functionality.


- If you have better structure idea contact to the provided email.

## Security

As `goxlock` is made for security , its base security features must not be altered to keep the application in good state.

Any design that do these cant be added to the application:

- Weaken the based action mecahnics.
- Storing data in any form without any need (plaintext must be used).
- Reduces integrity of the password and GCM.
- Bypass base authentication.
- Connects to net for an absurd reason.

And for security reasons:

- Do not publickly publish the security vulnerability , instead contact the provided email
- Those can posses risk onto the user and give advantage to the attacker.

## Pull Request

Before publishing your help please ensure that:

- The program runs successfuly and build a safe binary.
- No bugs in the written program.
- Passes all doctor test (Write your own test in the `doctor.Doctor.Start` to see it fits with other or not , try connecting it).

## Feature Request

If you got an idea , you can discuss about it before implementing it or sending it.
The email id is provided at the starting.

# TASKS

## Easy Tasks

- Fix typos in the the project.
- Add useful yet structural comments to them project.
- Any feature that is not documented , printed or used.

## Medium Tasks

- Testing current standards of code and executable and reporting which aspect it misses.
- Improve `--doctor` command (this is made for all the function that do writes onto the disk).
- Make `test` folder for `goxlock` to run safe and test non disk writting functions.
- Find Bugs and report (if security vulnerability , dont publish public , email to the givem gmail.)
- Give good looks to the cli interface.
- Provide useful messages in the code (if not provided).
- Optimizing the code for the least overhead.
- Help in making TUI for the app (using charmbracelet/bubbletea is appreciated).

## Hard Tasks

- Help making the binary for linux machines (providing function and details).

# BUGS

- Temp file early cleaneup can mess up scheduled tasks in Temp folder.