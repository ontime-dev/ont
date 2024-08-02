<div align="center" style="margin-bottom: 1em;">

Ontime
---

<img src="./docs/assets/images/ontimelogogreen.png" alt="ontime Logo" width=400px height=400px></img>
</div>

Ont (short for ontime) is a cronjob-like command line tool featuring simpler commands and more robust job management. With ont, you can store jobs in a database, stop them, and restart them at any time.

This project is in the beta phase, and using it in production environments is not recommended at the moment. However, we would appreciate it if you could test it and provide feedback.

## How is it different from cronjob and systemd timers?
- Easier syntax. You just need to pass the correct value to the option, pass the script path, and press enter.
- It is a matter of just one command line which can be executed from anywhere on the cluster to submit a job. 
- You can easily stop the job and start it again anytime later.
- You can check the status of the last execution of a job (in progress).
- You can run a job on a remote server and manage all the jobs from one single node and from any node connected to ontd in the cluster.
- Easy to manage and debug.

## Get started:
### For administrators:
> Note: The admin need to be either root user or sudo user.
1. Install ont on a server (check the installation guide).
2. Run the script initOnt.sh to create the DB, choose the DB password, and initialze ontd.
```
$ wget <script path>
$ chmod +x init_ont.sh
$ ./initOnt.sh
```
3. Modify the default values in the file ```/etc/ont/ont.conf``` as per the environment requirements.
4. Start the ont daemon as follows:
```
$ ont daemon
```
The ont daemon can also be started using systemd:
```
$ systemctl status ontd.service
```
If the file ontd.service doesn't exist, please create it and start the service. You can find a sample in ontd.service file.

5. To test the installation you can check the version and the help function.
```
$ ont -v
$ ont -h
```


### For users:
Using ```ont``` as a user is pretty simple, you just need to run one of the subcommands and pass the arguments.
- To submit a job:
  ```
  $ ont run --every 1d --from now /path/to/script.sh
  ```
  This will execute the script ```script.sh``` every day (1d) starting from now (now). Using ```now``` will execute the job immediately once, and then the next execution will be the following day at the same time. You can also pass a specific time or date. More examples can be found in the help output of the command.
  ```
  $ ont run -h
  ```

- To list all the jobs:
  ```
  $ ont list
  ```

- To start/stop/remove a job. You can use one of the subcommands ```start```, ```stop```, ```remove``` and pass the jobid.
  ``` 
  $ ont start <jobID>
  ```

  You can also run a job on a remote server by passing the name or the IP address of the remote server to the ```run``` subcommand.
   ```
   $ ont run -e 1d -f now -n node001 /path/to/script.sh
   ```
   This will execute the script ```script.sh``` on the node ```node001``` every day starting from now.

---
   ## Installation





