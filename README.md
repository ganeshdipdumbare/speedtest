# gospeedtest
gospeedtest is a simple command line tool to test the internet speed.

![](https://github.com/ganeshdipdumbare/gospeedtest/blob/main/demo.gif)

## Description
- The gospeedtest is a command line tool to test internet speed. It check the speed from `fast.com` and shows it on TUI.

- It uses [playwright-go](https://github.com/playwright-community/playwright-go) to scrape the website and get the download and upload speed. For showing it on beautiful TUI it uses [bubbletea](https://github.com/charmbracelet/bubbletea).

*NOTE*: The required dependencies for headless browser is fetched when the tool is run for the first time.

## Installation
- Requirements:-
    - Go 1.19
- To install the tool, just run the following command
    ```bash
    go install github.com/ganeshdipdumbare/gospeedtest@v0.1.2
    ```
## Improvements
- Add unit tests 
- Add other websites for speed check and allow user to select from the given options
