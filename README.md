# width-checker

`width-checker` is a cli tool to check if a page's width exceeds the device's width.  
This supports `firefox` and `chromium` by default. (powered by playwright)

## Usage

```
NAME:
   width-checker - A cli tool to check if a page's width exceeds the device's width.

USAGE:
   main [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --max value       Max is the maximum width of the device to be checked. (default: 1440)
   --min value       Min is the maximum width of the device to be checked. (default: 320)
   --interval value  Interval is the value of the change in width of the device. (default: 50)
   --url value       Url is the url of the page to be checked.
   --output, -o      (default: false)
   --help, -h        show help (default: false)

```

## Demo
