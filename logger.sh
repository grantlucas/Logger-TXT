#!/bin/bash

version()
{
  sed -e 's/^//' <<EndVersion
Logger.txt
Version 1.1
Author:  Grant Lucas (www.grantlucas.com)
Last updated:  31/01/2011
Release date:  26/07/2010
EndVersion
  exit 1
}

usage()
{
  sed -e 's/^//' <<EndUsage
Usage: logger.sh [-hV] [-t type] [-p project] [-d count] [-s] text
Try 'logger.sh -h' for more information.
EndUsage
  exit 1
}

help()
{
  sed -e 's/^//' <<EndHelp
Usage: logger.sh [-hV] [-t type] [-p project] [-d count] [-s] text

With no options or input, logger.sh outputs the last 10 lines of the log.

Options:
  -t TYPE
    The type classification that the log event belongs to. example: work, school etc.
  -p PROJECT
    The project that the log event belongs to. This helps group log events together which might belong to the same type or which my not belong to a type at all.
  -d COUNT
    The number of lines to show when output the tail of the log. Defaults to 10.
  -s text
    Searches the log file for the given text and displays those entries
  -h
    Help Text.
  -V
    Show version information and credits.
EndHelp
  exit 1
}

check_log_file()
{
  if [ -e $LOG_PATH ]; then
    if [ ! -w $LOG_PATH ]; then
      echo "$app: Log file not writeable"
      exit 1
    fi
  else
    # create log file if it does not exist
    echo "$app: Creating log file"
    `touch $LOG_PATH`
    `chmod +w $LOG_PATH`
    if [ -e $LOG_PATH ]; then
      echo "$app: Log file successfully created"
    else
      echo "$app: Log file couldn't be created"
      exit 1
    fi
  fi

  if [ ! -r $LOG_PATH ]; then
    echo "$app: Log file is not readable"
    exit 1
  fi 
}

search_log()
{
  #search the log for the serach term
  check_log_file
  #grep through file looking for the lines which have this
  results=`sed = "$LOG_PATH" | grep -i $SEARCH`
  echo -e "$results"
  exit 0
}

# defaults if not yet defined
dir=`dirname $0`
LOG_PATH=$dir"/log.txt"
LOG_TYPE=${LOG_TYPE:-''}
LOG_DISPLAY_COUNT=${LOG_DISPLAY_COUNT:-10}
LOG_PROJ=${LOG_PROJ:-''}

now=`date '+%d/%m/%y %H:%M'`
app="Log"

# process options
while getopts t:d:p:s:Vh o
do  case "$o" in
  s) SEARCH=$OPTARG
    search_log
  ;;
  t) LOG_TYPE=`echo "$OPTARG" | tr "[:lower:]" "[:upper:]"`;;
  d) LOG_DISPLAY_COUNT=$OPTARG;;
  p) LOG_PROJ=`echo "$OPTARG" | tr "[:lower:]" "[:upper:]"`;;
  h) help;;
  V) version;;
  [?]) usage;;
  esac
done
# shift the option values out
shift $(($OPTIND - 1))

#The remaining text is the log text.
#FIXME: Escape log text of special characters which will mess up the insert. Mainly $'s
#log_text=`echo "$*" | sed -e 's/\$/\\\$/g'`
#echo $log_text

#exit 1


#take the input and add to file
if [ ! -z "$1" ]; then
  #add to log file
  check_log_file

  if [ ! -z $LOG_TYPE ]; then
    sep=" - "
    ltype=" under the type $LOG_TYPE"
    LOG_TYPE="$LOG_TYPE"
  fi

  if [ ! -z $LOG_PROJ ]; then
    sep=" - "
    proj=" in the project $LOG_PROJ"
    LOG_PROJ="($LOG_PROJ)"
  fi

  #there is a proj but no type
  if [ -z $LOG_TYPE ] && [ ! -z $LOG_PROJ ]; then
    category="$LOG_PROJ$sep"
  fi

  #there is a type but no proj
  if [ ! -z $LOG_TYPE ] && [ -z $LOG_PROJ ]; then
    category="$LOG_TYPE$sep"
  fi

  #there is both
  if [ ! -z $LOG_TYPE ] && [ ! -z $LOG_PROJ ]; then
    category="$LOG_TYPE $LOG_PROJ$sep"
  fi

  #add text to file
  echo "$now - $category$*" >> "$LOG_PATH" 
  #output that the event was logged
  echo "$app: $* logged$ltype$proj"
else
  #no options so print out line by line log file
  check_log_file
  tail -r -n $LOG_DISPLAY_COUNT $LOG_PATH
fi
