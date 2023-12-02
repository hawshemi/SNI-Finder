# SNI Finder

This script will scan all domains that have `TLS 1.3` and `h2` enabled on your VPS IP address range. These domains are useful for SNI in a variety of configurations and tests.

When you begin the scan, two files are created: `results.txt` contains the output log, while `domains.txt` only contains the domain names.

It is Recommended to run this scanner locally _(with you residential internet)_. It may cause VPS to be flagged if you run scanner in the cloud.


## Run

### Linux:

1.
    ```
    wget "https://github.com/hawshemi/SNI-Finder/releases/latest/download/SNI-Finder-linux-amd64" -O SNI-Finder && chmod +x SNI-Finder
    ```
2. 
    ```
    ./SNI-Finder -addr ip
    ```

### Windows:

1. Download from [Releases](https://github.com/hawshemi/SNI-Finder/releases/latest).
2. Open `CMD` or `Powershell` in the directory.
3.
    ```
    .\SNI-Finder-windows-amd64.exe -addr ip
    ```

#### Replace `ip` with your VPS IP Address.


## Build

### Prerequisites

#### Install `wget`:
```
sudo apt install -y wget
```

#### First run this script to install `Go` & other dependencies:
```
wget "https://raw.githubusercontent.com/hawshemi/SNI-Finder/main/install-go.sh" -O install-go.sh && chmod +x install-go.sh && bash install-go.sh
```
##### Reboot is recommended.


#### Then:

#### 1. Clone the repository
```
git clone https://github.com/hawshemi/SNI-Finder.git 
```

#### 2. Navigate into the repository directory
```
cd SNI-Finder 
```

#### 3. Download and install `logrus` package
```
go get github.com/sirupsen/logrus
```

#### 4. Build
```
CGO_ENABLED=0 go build
```
