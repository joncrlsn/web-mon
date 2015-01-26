package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	templ "text/template"
)

var dumpThreadsScript = `#!/bin/bash
# Dumps Java threads a certain number of times
ssh -o StrictHostKeyChecking=no {{.host}}.c42 'bash -s' <<-END
#!/bin/bash
COUNT={{.dumpCount}}
FILE="/home/{{.cpuser}}/logs/threads.\$(date +%Y-%m-%d_%H%M%S).{{.host}}"
PID=\$(ps aux | grep -P '(central|blue).*java' | grep -v grep | grep -v flock | egrep -v 'su (central|blue)' | awk '{print \$2}')
#echo "\$FILE"
#echo "\$PID"
for (( c=1; c<=COUNT; c++ )) ; do
    sudo su {{.cpuser}} -- -c "touch \${FILE}; jstack -l \$PID >> \${FILE}"
    echo "Threads dumped... to \$FILE.  Sleeping for {{.intervalSeconds}} seconds..."
    sleep {{.intervalSeconds}}
done
echo done
END
`

// dumpJavaThreds dumps the Java threads on the given host for the given number of times
// This will need modification for monitoring a Windows-based resource because it creates a 
// bash script that is executed remotely.
func dumpJavaThreads(host, user string, dumpCount int, intervalSeconds int) error {

	// Save script file
	filename := host + "_dumpThreadScript.sh"
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	fmt.Println(" Writing script to file : " + filename)

	t := templ.New("Dump threads script")
	templ.Must(t.Parse(dumpThreadsScript))
	ctx := map[string]string{"cpuser": user,
		"dumpCount":       strconv.Itoa(dumpCount),
		"intervalSeconds": strconv.Itoa(intervalSeconds),
		"host":            host}
	err = t.Execute(file, ctx)
	if err != nil {
		return err
	}
	file.Close()

	cmdStr := "./" + host + "_dumpThreadScript.sh"
	fmt.Printf("Executing: %s \n", cmdStr)
	sshCmd := exec.Command(cmdStr)
	bytes, err := sshCmd.CombinedOutput()
	fmt.Printf("Output:\n %s \n", string(bytes))
	if err != nil {
		return err
	}

	return nil
}
