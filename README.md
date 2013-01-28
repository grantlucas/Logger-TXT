Logger TXT - Quick command line logging
=======================================

Logger TXT is a small, shell based tool to log activities throughout the day to a simple, portable text file, along with the date/time. Options are available to log a specific entry under a type and project. All entries are stored in a simple TXT file. Whether you track purchases, what you ate that day, progress on projects at work or all of the above and more, you will always have a simple, solid way of storing that information and a script that gets out of your way to get it there.

Installation
============

Location of log script
----------------------

To install loggerTXT simply copy logger.sh to where you would like it to be stored on your computer. If you have multiple computers logger works really well within a folder in Dropbox. For example you could keep logger.sh and your log file in ~/Dropbox/log/ and it will automatically be synced between computers.

Location of log file
--------------------

There are two options for defining where your log file is to be saved.

1) Setting an environment variable

  - In ~/.profile add:

    export LOGGERTXT\_PATH=~/Dropbox/log/log.txt

  - Make sure you include the name of your log file. This allows you to set it to a hidden file if you desire.

2) Default action if no environment variable set

  - If no LOGGERTXT\_PATH is set, a log.txt file will be created in the folder where logger.sh is located.


Quick Command Line Access
=========================
In ~/.profile add:

    alias l="/path/to/script/logger.sh"

Example Input
=============

Without Alias
-------------

    ./logger.sh -t personal -p project "This is a log note with a type and project"

With Alias
----------

    l -t personal -p project "This is a log note with a type and project"

Example output in log.txt
=========================

    31/01/11 13:30 - PERSONAL (PROJECT) - This is a log note with a type and project
    31/01/11 13:35 - PERSONAL - This is a log not with just a type
    31/01/11 13:40 - (PROJECT) - This is a log not with just a project
    31/01/11 13:45 - This is just a general event which doesn't belong to anything

Main Goals
==========

The main goal of this project was to create a simple logging tool which could be accessed quickly from within the command line environment. By storing all data in a TXT file, you're not locked into always using this tool or limited to only viewing log events with this script. The data portability that a TXT file offers between tools, operating systems and environments is crucial to having a smooth workflow that is extremely dependable.

What do you use it for anyways?!?
=================================

Over time the act of logging will become habitual. Over the course of a day you may log any of the following and anything else you deem important.

- Progress of tasks related to work and/or specific projects
  - Extremely handy when it comes to filling in hours with an employer as you can easily look up what projects were worked on, on that Tuesday two weeks ago.
- Progress of personal tasks or projects
  - Progress logging is the main use of this tool
- Purchases made
  - Extremely useful when the credit card bill comes with cryptic names of companies.
- Log important events or anything where the time that it happened is important.
  - Had an important conversation with someone? Log that you had it so you can also know when it exactly happened.
- Log anything!

Tools to make it simpler and easier
===================================

OS X
----

- [iTerm2](http://www.iterm2.com/)
  - **Recommended**
  - More advanced program than the stock Terminal.app
  - Comes with a built in quick view window so you can quickly bring up a terminal window to user LoggerTXT
- [Visor](http://visor.binaryage.com/)
  - Terminal.app symbol plug-in which keeps a terminal window always available with a quick keyboard shortcut as per the game Doom and ~
