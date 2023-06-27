

# Host-CLI
Host-CLI is a Ads blocker and websites blocker CLI based tool written in golang.

## Installation ( Debian/Ubuntu )

* [Download](https://github.com/BiswajitThakur/host-cli/releases) the Debian package & Execute the following command
```bash
sudo dpkg -i host-cli_*.deb
sudo host-cli --version
```
## Installation ( Linux )
```bash
sudo apt install golang
git clone https://github.com/BiswajitThakur/host-cli
cd host-cli/
go build .
sudo ./host-cli --version
```
## Installation ( Windows )
* [Download](https://github.com/BiswajitThakur/host-cli/releases) the exe file.
* Open terminal as administrator.
* Go to the location where you downloaded the exe file.
* Execute the following command.
```bash
host-cli --version
```
* If the above command print version, please read Documentation

## Documentation

### - Block Ads -
```bash
sudo host-cli --block
```

### - Unblock Ads -
```bash
sudo host-cli --unblock
```

### - Block Single Website -
```bash
sudo host-cli --block <host_name>
```
For example : `sudo host-cli --block google.com`

### - Block Multiple Website -
```bash
sudo host-cli --block <host_name1> <host_name2> <host_name3> ...
```
For example : `sudo host-cli --block google.com facebook.com`

### - Unblock Single Website -
```bash
sudo host-cli --unblock <host_name>
```
For Example : `sudo host-cli --unblock google.com`

### - Unblock Multiple Website -
```bash
sudo host-cli --unblock <host_name1> <host_name2> <host_name3> ...
```
For Example : `sudo host-cli --unblock google.com facebook.com`

### - Update hosts sources -
`sudo host-cli --updateSourceList` or `sudo host-cli --upsl`

### - Create http server for using GUI mode on browser
`sodo host-cli --http 3999` or `sodo host-cli --http 127.0.0.1:3999`
* Now open [http://127.0.0.1:3999](http://127.0.0.1:3999)
