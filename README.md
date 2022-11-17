# trollgo
Trollgo periodically blocks mouse movements, modifies keyboard inputs and swaps mouse buttons
# Build from source
Clone this repository and move into the directory

    git clone https://github.com/sne4ker/trollgo.git
    cd trollgo

Build with GUI

    GOOS=windows go build -ldflags "-s -w" -o trollgo.exe

Build without GUI

    GOOS=windows go build -ldflags "-H=windowsgui -s -w" -o trollgoNogui.exe

# TODO
- Move mouse randomly
- (Play sound)
