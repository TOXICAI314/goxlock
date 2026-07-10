# UPDATE INFO

Following are the changes made in the `v1.1` update:

## Fixed Bugs 

- Missing File stats in the end -> Because the file was already deleted in the first place by `--del-original`.

- The Above bug made the app panic in giving status while enabling `--del-original` and `--stats` in the same time.

- The double click file mechanics will now delete the locked file on spot after unlcoking.

- `unsafe` use for time out prompt bypass. (As `unsafe` was meant to bypass these.)


## Design Improvement

- Made the error message better -> Classified in two factors : FunctionCacelError (for self killing) & FucntionFailError (killed by an error).

- Better Error messages are given.

- Made folder imoprovements.

- Added `CONTIBUTING.md` for a better display.

## Delayed Updates

Due to developemennt issues and lack of support , these features have been dropped to future updates.

- UI UPDATE (because `bubbletea` is taking too long to understand and my exams are coming.)

## Focused Updates

- Linux function and binary support.