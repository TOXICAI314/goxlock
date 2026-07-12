# UPDATE INFO

Following are the changes made in the `v1.0.1-Beta` update:

## Fixed Bugs 

- Bug : Missing File stats in the end -> Because the file was already deleted in the first place by `--del-original`.

- Bug : The double click file mechanics will dont delete the locked file on spot after unlocking.

- Bug : `unsafe` use for time out prompt bypass. (As `unsafe` was meant to bypass these.).

- Buggy `--update` command locking its own way because of the file not closed.

- Version Check instability.

- The `--header` bug where it tries to give the stringed version in form of int (via binary manipulation).

- `os.Size` for a folder like object will yield a number like `4096`.

- `Mutex` creation has a hidden bug of wasting the mutex when the a hard coded condition is met.

- `--change-password` in `PerformAction` returns when not needed

- `ShowSensitiveData` write the data into the Buffer which then cant be controlled by the application to clear.So the feature to see the data is dropped for the future GUI Update.

- The slip made by the folder by escaping the current directory.

- Updater not working because of the self app is not exiting.

- The client will never end even if no response comes to it for a long period.

## Design Improvement

- Made the error message better -> Classified in two factors : FunctionCacelError (for self killing) & FucntionFailError (killed by an error).

- Better Error messages are given.

- Made folder imoprovements.

- Added `CONTIBUTING.md` for a better display.

- Made the `Version` as Spearated operating constants , like : `Major`,`Minor`,`Patch`,`Release`.

- Made the `actions` as `enums`

- Validation logic becomes more optimised and confesticated.

- Made the `Appdata` -> `Config` pointing. Where it can be changed as needed in the future.

- `logger` now logs if the code mitsakely returns in the switch conditions.

- `ReplaceFolderwithGlock` cant delete the volume name even if it is prompted to by `--unsafe`

- Added a delay to the Updater so that it gets a valid time to react.

## Delayed Updates

Due to developemennt issues and lack of support , these features have been dropped to future updates.

- UI UPDATE (because `bubbletea` is taking too long to understand and my exams are coming.)

## Focused Updates

- Linux function and binary support.