<h1 align="center">Camplified</h1>

<div align="center">

<img src="icon.ico"></img>

[[Download Latest Release](https://github.com/HarukaKinen/Camplified-go/releases/latest)]

Yet another osu! leaderboard camping program but simple implementation in Golang

![](https://img.shields.io/github/go-mod/go-version/HarukaKinen/Camplified-go/master?style=for-the-badge)
[![](https://img.shields.io/github/v/release/HarukaKinen/Camplified-go?style=for-the-badge)](https://github.com/HarukaKinen/Camplified-go/releases/latest)
![](https://img.shields.io/github/license/HarukaKinen/Camplified-go?style=for-the-badge)

</div>

## Troubleshooting

See some weird code like ``[1;1H``, ``[2;4H`` in your terminal?

Your Windows terminal seem doesn't enable the support of **ANSI escape sequences**.

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
