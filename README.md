# jinglepings-plot
hacky project to plot images on https://jinglepings.com/

## Usage

1. Download/Clone the repository
2. run `go get .../.` to install dependencies
3. run `go build -o jinglepings`
4. run `sudo ./jinglepings path/to/image.png`


(requires root to send ICMP pings in raw socket mode)

## Run continously (ish) 

Run every 10 seconds: `watch -n10 sudo ./jinglepings pickle.png`


