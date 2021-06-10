# Resources 

* [Understanding Systemd Units and Unit Files](https://www.digitalocean.com/community/tutorials/understanding-systemd-units-and-unit-files)
* [Digital Ocean Systemd guide](https://www.digitalocean.com/community/tutorials/understanding-systemd-units-and-unit-files)
* [Running go binary as systemd service](https://fabianlee.org/2017/05/21/golang-running-a-go-binary-as-a-systemd-service-on-ubuntu-16-04/)


# Running this chess engine as a systemd service on Ubuntu 

1. Build the executable
    ```
    $ go build -o glee cmd/glee/main.go
    ```
1. Move the binary to a folder in your path, e.g.
    ```
    $ sudo cp glee /usr/local/bin/
    ```
1. Update `glee.service` ConditionPathExists, User, Group, WorkingDirectory, and ExecStart as needed.
1. Move `glee.service` to `/lib/systemd/system/glee.service`
    ```
    $ sudo cp config/systemctl/glee.service /lib/systemd/system/
    ```
1. Update the file permissions
    ```
    $ sudo chmod 755 /lib/systemd/system/glee.service
    ```
1. Start the service 
    ```
    $ sudo systemctl glee start
    $ // if this fails with 'unknown operation' try `sudo systemctl start glee`
    ```
1. Monitor the service ouput 
    ```
    $ journalctl -f -u glee
    ```
1. Enable the service to start at boot 
    ```
    $ sudo systemctl enable glee
    ```
