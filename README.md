# zerol-enrollment-cleaner

`zero-l-enrollment-cleaner` is a command line tool that can be used to remove enrollees from HLS Zero-L courses in the ExEd Canvas instance.

## Getting started

1. View the details for the latest release on [GitHub](https://github.com/Harvard-University-iCommons/zerol-enrollment-cleaner/releases/latest), and download the version for your platform (Mac or Windows).
2. See the help for the `zerol-enrollment-cleaner` command line tool by running `zerol-enrollment-cleaner --help` (or `zerol-enrollment-cleaner.exe -h` on Windows).

## Usage

```
Usage of zerol-enrollment-cleaner:
  -account_id int
        The Canvas account ID to use (default 139, HLS Online)
  -course_id int
        The Canvas course ID to use
  -file string
        File to read - must contain a list of email addresses, one per line
  -host string
        The Canvas host to connect to (default "exed.canvas.harvard.edu")
  -token string
        The API token to use
```

So, for example:

```
> zerol-enrollment-cleaner -account_id 139 -course_id 1234 -file my_enrollments.txt -token my_token
```
or (Windows):
```
c:\\> zerol-enrollment-cleaner.exe -account_id 139 -course_id 1234 -file my_enrollments.txt -token my_token
```

# Developer information

This tool is written in [Go](https://golang.org), primarily to make it easy for us to provide self-contained binaries compiled for various platforms.

During development, you can run the tool from the command line with `go run .`.

To see which platforms you can compile for, run `go tool dist list`.

To compile for a specific platform, run:
`GOOS=windows go build` for Windows, or
`GOOS=darwin go build` for Mac.
