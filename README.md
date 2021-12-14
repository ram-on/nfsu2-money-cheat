# nfsu2-money-cheat
Money cheat for "Need For Speed Underground 2" -- allows you to edit/change money within your NFSU2 save file.

## Download EXE

Link:  https://github.com/ram-on/nfsu2-money-cheat/releases/download/v1.0/nfsu2-money-cheat.exe

## Examples

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

## Compatibility

Program can be compiled and executed on any platform supported by the GO language.  Tested on Windows and MacOS.
