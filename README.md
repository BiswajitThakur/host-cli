
# Host-CLI
Host-CLI is a Ads blocker and websites blocker CLI based tool written in golang.

## Installation ( Debian/Ubuntu )

* Clone this repository -

```bash
git clone https://github.com/BiswajitThakur/host-cli

```
* Now go to cloned directory -

```bash
cd host-cli/
```

* Run setup.sh
```bash
sudo bash setup.sh
```

* Now you will see `host-cli` executable file. Execute this file using `sudo` command -
```bash
sudo ./host-cli --version
```

## Installation ( Termux - For Rooted Devices)
* Install dependencies
```bash
pkg update && pkg install tsu golang
```
* Clone this repository -

```bash
git clone https://github.com/BiswajitThakur/host-cli

```

* Execute these following commands -

```bash
cd host-cli/
```
```bash
bash setup.sh
```

* Now you will see `host-cli` executable file. Execute this file using `sudo` command -
```bash
sudo ./host-cli --version
```

## Documentation

### - Block Ads -
```bash
sudo ./host-cli --block
```

### - Unblock Ads -
```bash
sudo ./host-cli --unblock
```

### - Block Single Website -
```bash
sudo ./host-cli --block <host_name>
```
For example : `sudo ./host-cli --block google.com`

### - Block Multiple Website -
```bash
sudo ./host-cli --block <host_name1> <host_name2> <host_name3> ...
```
For example : `sudo ./host-cli --block google.com facebook.com`

### - Unblock Single Website -
```bash
sudo ./host-cli --unblock <host_name>
```
For Example : `sudo ./host-cli --unblock google.com`

### - Unblock Multiple Website -
```bash
sudo ./host-cli --unblock <host_name1> <host_name2> <host_name3> ...
```
For Example : `sudo ./host-cli --unblock google.com facebook.com`
