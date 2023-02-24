# speedtest
speedtest is a simple command line tool to test the internet speed.

![](https://github.com/ganeshdipdumbare/speedtest/blob/main/demo.gif)

## Description
- The speedtest is a command line tool to test internet speed. It check the speed from `fast.com` and shows it on TUI.

- It uses [playwright-go](https://github.com/playwright-community/playwright-go) to scrape the website and get the download and upload speed. For showing it on beautiful TUI it uses [bubbletea](https://github.com/charmbracelet/bubbletea).

*NOTE*: The required dependencies for headless browser is fetched when the tool is run for the first time.

## Installation
- Requirements:-
    - Without Golang 
        ```bash
        curl -sf https://gobinaries.com/ganeshdipdumbare/speedtest | sh
        ```
    - With Golang - Go 1.19+
        - To install the tool, just run the following command
        ```bash
        go install github.com/ganeshdipdumbare/speedtest@v0.1.4
        ```
## Improvements
- Add unit tests 
- Add other websites for speed check and allow user to select from the given options
