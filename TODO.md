Todo
====

Searching
---------
- Add case sensitive search as -S option
  - Update help text for case sensitive/insensitive searching
- Add ability to search specifically by type and/or project 
  - examples
    - Just show all "WORK" related entries
    - Just show all "(PROJECT)" related entries
    - Show all work items from a specific project

Improvements
------------
- If no log file exists and the path isn't set in an environment variable, but the log file in a folder in the user's home folder instead of where ever the log script is stored.
  - For example, if the user stores logger.sh in /usr/local/bin and there is _no_ environment variable defining the path to store the file, it will be created in /usr/local/bin which is baaaaaaaddddd.
  - Must fix as soon as possible
- Ability to create a status report
  - Still formulating how this will work and what the result will be
- Web App instead of native mobile app
  - Ties into dropbox
  - Do everything that I'd want to do with a native app, but on the web
  - Write in zend?
- Add time limit for deleting the last line from the file
  - Make is so that this must be triggered within a time limit from the time the last line was entered. This will help prevent accidental deletions later in time

Bugs
----
- $'s are not properly escaped
- !'s are not properly escaped
