<h1 align="center">Camplified</h1>

<div align="center">

<img src="icon.ico"></img>

[[Download Latest Release](https://github.com/HarukaKinen/Camplified-go/releases/latest)]

Yet another osu! leaderboard camping program but simple implementation in Golang

[![](https://img.shields.io/github/go-mod/go-version/HarukaKinen/Camplified-go/master?style=for-the-badge)](https://go.dev/)
[![](https://img.shields.io/github/v/release/HarukaKinen/Camplified-go?style=for-the-badge)](https://github.com/HarukaKinen/Camplified-go/releases/latest)
[![](https://img.shields.io/github/license/HarukaKinen/Camplified-go?style=for-the-badge)](https://github.com/HarukaKinen/Camplified-go/blob/main/LICENSE)

[![DeepSource](https://app.deepsource.com/gh/HarukaKinen/Camplified-go.svg/?label=active+issues&show_trend=true&token=8SZvFbqeextkNRaPHU2ep_bV)](https://app.deepsource.com/gh/HarukaKinen/Camplified-go/?ref=repository-badge)

</div>

## How does it work?

An endless loop of API request to check map's status, if the map is ranked, it will play a sound to notify you.

### Sound Effect

Campilfied's sound is based on Windows' ``kernel32.dll``'s [``Beep``](https://learn.microsoft.com/en-us/windows/win32/api/utilapiset/nf-utilapiset-beep) function. This is a part of Windows API. If your Windows doesn't have ``kernel32.dll``, you **can't run** Campilfied. This should't be happen in most cases.

## Troubleshooting

See some weird code like ``[1;1H``, ``[2;4H`` in your terminal?

Your Windows seem doesn't enable the support of **ANSI escape sequences**.

But you have multiple ways to fix it.

### Enable support of ANSI Terminal Control

Run **PowerShell** as ***administrator*** and execute the following command:

```powershell
Set-ItemProperty HKCU:\Console VirtualTerminalLevel -Type DWORD 1
```

This command will change the registry value of ``VirtualTerminalLevel`` in ``HKEY_CURRENT_USER\Console\`` to ``1``. Aka enable the support of ANSI Terminal Control. Also change to ``0`` means disable.

Learn more in [Windows console with ANSI colors handling - Stack Exchange](https://superuser.com/a/1300251/1803960)

### Run Camplified in Mircosoft Terminal

Feel restless to edit the registry? Try [Microsoft Terminal](https://github.com/microsoft/terminal)

- Get it from Microsoft Store: [Windows Terminal](https://aka.ms/terminal)
- Download the latest release from [GitHub](https://github.com/microsoft/terminal/releases/latest)
- You can also learn more about how to installing and running Windows Terminal in their [Github repository](https://github.com/microsoft/terminal#installing-and-running-windows-terminal)

After installing, you have multiple ways to run Camplified in Microsoft Terminal. We take a simple way here.

1. Going to the directory of ``Camplified.exe``.
2. Right click empty space and click ``Open in Windows Terminal``, or type ``wt`` in the address bar of File Explorer to open Microsoft Terminal in current directory.
3. Then type ``.\Camplified.exe`` to run Camplified.
4. You're done!
