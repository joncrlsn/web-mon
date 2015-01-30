#!/bin/bash
#
# This is an example script one could write to do some processing when an alert is triggered.
# Name it whatever you want and add it to the shellCommand variable in your config file. 
# 
# This script remotely dumps Java threads a number of times with an interval begin thread dumps
#
HOST=$1
PROCESS_OWNER=$2
ssh -o StrictHostKeyChecking=no $HOST 'bash -s' <<-END
#!/bin/bash
COUNT=10
INTERVAL=8
# Write the thread dumps to a particular location
FILE="/home/${PROCESS_OWNER}/logs/threads.$(date +%Y-%m-%d_%H%M%S).$HOST"
PID=$(ps aux | grep -P '(central|blue).*java' | grep -v grep | grep -v flock | egrep -v 'su (central|blue)' | awk '{print $2}')
echo "$FILE"
echo "$PID"
for (( c=1; c<=COUNT; c++ )) ; do
    sudo su $PROCESS_OWNER -- -c "touch ${FILE}; jstack -l $PID >> ${FILE}"
    echo "Threads dumped... to $FILE.  Sleeping for $INTERVAL seconds..."
    sleep $INTERVAL
done
echo done
END
