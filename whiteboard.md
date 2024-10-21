
▀██▀▀█▄                    ▄
 ██   ██  ▄▄▄▄     ▄▄▄▄  ▄██▄
 ██▄▄▄█▀ ▀▀ ▄██  ▄█   ▀▀  ██
 ██      ▄█▀ ██  ██       ██
▄██▄     ▀█▄▄▀█▀  ▀█▄▄▄▀  ▀█▄▀

Like Apache Airflow, but simpler and... I don't know, just simpler.

But hey, it's still fair - sometimes I just want to schedule sync.
~~*systemd/Timers*~~

> As a user who is running shell scripts,
I want to run them with intervals in the background.
I also want to be sure that session - background activity - will continue after startup.
I care about preciseness and ease of use, sovereignty and interoperability.

Alright, so user has to choose a folder, that will contain scripts they want to run, schedule or delay.
~~These scripts will have to have a specific format(see *RFC: File-based metadata*),
like a header with metadata, and then the script itself.~~
Scripts will have a TOML file with metadata.

~~**RFC: File-based metadata**~~
> ~~*Instead of forcing user to edit their scripts writing metadata, derive it from the filename.*~~
>
> ~~*Pain-points: cron.*~~

**Metadata**

- `name` (required):
  - The identifier of the script.
- `file` (optional):
  - The path to the script file. If not provided, file will be matched with name.
- `schedule` (optional):
  - The schedule in cron format.
- `at` (optional):
  - The time to run the script in ISO 8601 format. If both schedule and at are provided, `at` will take precedence.
- `description` (optional):
  - A description of the script.
- `active` (optional):
  - Whether the script is active or not. Default is `true`.

10/19/24: updated from example folder

MVP:
- [x] config: set scripts directory
- [x] create a structure for the metadata
- [x] parse the scripts for metadata
- [ ] example: run beeps at different frequencies(a4, g4, f4) and intervals(1s, 2s, 3s)
- [ ] run the scripts with the correct interpreter (shebang)

// how to design Demon reading/writing from/to scripts.toml at runtime?
- [ ] watch for changes in the scripts directory and update running Demon instance

- [ ] test: Demon will run the scripts according to the schedule
- [ ] test: Demon will log the output of the scripts
