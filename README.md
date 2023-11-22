# nfsu2-money-cheat
Money cheat for "Need For Speed Underground 2" -- allows you to edit/change money within your NFSU2 save file.

## Download EXE

Download Link:  https://github.com/ram-on/nfsu2-money-cheat/releases/download/v2.0/nfsu2-money-cheat.exe

## Examples

### Text User Interface

To launch the program in text user interface simply double click the `nfsu2-money-cheat.exe`.
User will then be asked to select the save file and to specify the money amount that will be
stored within the save file.

### View Money
To view the curreny money amount stored within your save file:

```
nfsu2-money-cheat.exe info -f "%LOCALAPPDATA%\NFS Underground 2\MySaveFile"
```

### Edit Money
To edit/change the money amount stored to "50000" within your save file:

```
nfsu2-money-cheat.exe edit -m 50000 -f "%LOCALAPPDATA%\NFS Underground 2\MySaveFile"
```

## Compile

Program can be compiled and executed on any platform supported by the GO language.
Tested on Windows and MacOS.

To compile on Windows via PowerShell:

```
$env:GOARCH=386; go build nfsu2-money-cheat.go; $env:GOARCH=$null
```

To compile on Windows via CMD.EXE:

```
set "GOARCH=386"  &&  go build nfsu2-money-cheat.go & set GOARCH=
```

To compile on Linux/macOS and produce a Win32 binary:

```
make build-win32
```

## License

Program licensed under [The Unlicense](https://github.com/ram-on/nfsu2-money-cheat/blob/main/LICENSE) license.
