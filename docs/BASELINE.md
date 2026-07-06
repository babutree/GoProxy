# Baseline Snapshot

Task: P0-01 Establish dev branch and baseline snapshot

Created at: 2026-07-07T00:41:24.8431443+08:00

## Baseline Commit

`49e99a526e6beaa518efd0392735a908902276ce`

## Toolchain

### Go

Command: `go version`

Result: failed. `go` is not available in the current PATH, so the required Go version check (`>= 1.25`) could not be verified.

Error output:

```text
go:
Line |
   2 |  go version
     |  ~~
     | The term 'go' is not recognized as a name of a cmdlet, function, script file, or executable program.
     | Check the spelling of the name, or if a path was included, verify that the path is correct and try again.
```

### CGO Compiler

Command: `gcc --version`

Result: failed. `gcc` is not available in the current PATH.

Additional MinGW/MSYS2 presence check:

Checked common gcc paths:

```text
C:\msys64\mingw64\bin\gcc.exe
C:\msys64\ucrt64\bin\gcc.exe
C:\mingw64\bin\gcc.exe
C:\MinGW\bin\gcc.exe
```

No gcc executable was found at those paths.

CGO compiler available: no

Error output:

```text
gcc:
Line |
   2 |  gcc --version
     |  ~~~
     | The term 'gcc' is not recognized as a name of a cmdlet, function, script file, or executable program.
     | Check the spelling of the name, or if a path was included, verify that the path is correct and try again.
```

## Build Verification

Command:

```powershell
$env:CGO_ENABLED='1'; go build ./...
```

Result: fail

Full error output:

```text
go:
Line |
   2 |  $env:CGO_ENABLED='1'; go build ./...
     |                        ~~
     | The term 'go' is not recognized as a name of a cmdlet, function, script file, or executable program.
     | Check the spelling of the name, or if a path was included, verify that the path is correct and try again.
```
