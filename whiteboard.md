
▀██▀▀█▄                    ▄
 ██   ██  ▄▄▄▄     ▄▄▄▄  ▄██▄
 ██▄▄▄█▀ ▀▀ ▄██  ▄█   ▀▀  ██
 ██      ▄█▀ ██  ██       ██
▄██▄     ▀█▄▄▀█▀  ▀█▄▄▄▀  ▀█▄▀

Like Apache Airflow, but simpler and... I don't know, just simpler.

But hey, it's still fair - sometimes I just want to schedule sync.
*systemd/Timers*

> As a user who is running shell scripts,
I want to run them with intervals in the background.
I also want to be sure that session - background activity - will continue after startup.
I care about preciseness and ease of use, sovereignty and interoperability.

Alright, so user has to choose a folder, that will contain scripts they want to run, schedule or delay.
These scripts will have to have a specific format(see *RFC: File-based metadata*),
like a header with metadata, and then the script itself.

**RFC: File-based metadata**
> *Instead of forcing user to edit their scripts writing metadata, derive it from the filename.*
>
> *Pain-points: cron.*


**Metadata.**
- `name: string` - an identifier,
- `schedule: CronExp | UnixTime`
  - if `CronExp` or *cron expression*:
    Schedules a script execution with cron - (execute script now or at the next scheduled time?)
  - if `UnixTime` or *unix timestamp*: Schedules a script execution at the fixed time



MVP:
- [x] config: set scripts directory
- [ ] parse the scripts for metadata
- [ ] create a structure for the metadata

- [ ] Demon will run the scripts according to the schedule
- [ ] Demon will log the output of the scripts

~~In case a script panics, we could try to fix errors, but it might take more than llvm repository to account for all cases.
However, we can define common cases:~~
- ~~file permissions, chmod~~
- ~~dependencies~~

---
10/16/24 8PM
hey, there is a tool called Task and it got 11.2k start on GitHub and written in Go, and
it is quite cool, because you can create a YAML file and describe all your tasks there.
Like "cmds" that is basically just shell script. Then you simply

zzz

10/18/24 1:17AM
hey, I created a config. I don't know what for yet. But I will soon discover.

10/18/24 1:20AM
I remembered - it was for setting scripts directory. See, `scripts` directory has a special meaning.
It needs scripts to be located in one folder, in return it will be easier to work with them
Maybe we could, instead of forcing user to edit their scripts writing metadata, just derive it from the filename.
Hmm, I don't know if cron would work this way. Maybe not the cron expressions - `*` is illegal -
but this: `sync-15m.sh`. And this: `tweetbot-1mo(7th).sh`
