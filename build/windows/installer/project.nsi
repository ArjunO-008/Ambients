Unicode true

####
## Ambients — Windows NSIS Installer
## Wails v2 project.nsi — customised for Ambients v1.0.0
##
## Build with:
##   wails build -platform windows/amd64 -clean -nsis
##
## Or manually (after a wails build populates wails_tools.nsh):
##   makensis -DARG_WAILS_AMD64_BINARY=..\..\bin\Ambients.exe
####

!include "wails_tools.nsh"

# ─── Version info ─────────────────────────────────────────────────────────────
VIProductVersion "1.0.0.0"
VIFileVersion    "1.0.0.0"

VIAddVersionKey "CompanyName"     "Arjun.O"
VIAddVersionKey "FileDescription" "Ambients Installer"
VIAddVersionKey "ProductVersion"  "1.0.0"
VIAddVersionKey "FileVersion"     "1.0.0"
VIAddVersionKey "LegalCopyright"  "MIT License, Copyright 2026 Arjun.O"
VIAddVersionKey "ProductName"     "Ambients"

# ─── HiDPI ────────────────────────────────────────────────────────────────────
ManifestDPIAware true

# ─── MUI setup ────────────────────────────────────────────────────────────────
!include "MUI.nsh"

!define MUI_ICON   "..\icon.ico"
!define MUI_UNICON "..\icon.ico"

# Welcome page image — optional, comment out if you don't have one
# !define MUI_WELCOMEFINISHPAGE_BITMAP "resources\welcome.bmp"

# Installer behaviour
!define MUI_FINISHPAGE_NOAUTOCLOSE
!define MUI_ABORTWARNING

# Finish page — offer to launch Ambients after install
!define MUI_FINISHPAGE_RUN          "$INSTDIR\Ambients.exe"
!define MUI_FINISHPAGE_RUN_TEXT     "Launch Ambients"

# Finish page — link to GitHub releases for future updates
!define MUI_FINISHPAGE_LINK         "Check for updates on GitHub"
!define MUI_FINISHPAGE_LINK_LOCATION "https://github.com/ArjunO-008/ambients/releases"

# ─── Pages ────────────────────────────────────────────────────────────────────
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_INSTFILES

!insertmacro MUI_LANGUAGE "English"

# ─── Installer metadata ───────────────────────────────────────────────────────
Name              "Ambients"
OutFile           "..\..\bin\Ambients-amd64-installer.exe"
InstallDir        "$LOCALAPPDATA\Ambients"
ShowInstDetails   show

# Install to user-local dir — no admin rights required
RequestExecutionLevel user

# ─── Architecture check ───────────────────────────────────────────────────────
Function .onInit
  !insertmacro wails.checkArchitecture
FunctionEnd

# ─── Install section ──────────────────────────────────────────────────────────
Section "Ambients" SecMain

  !insertmacro wails.setShellContext

  # Install WebView2 runtime if not already present
  !insertmacro wails.webview2runtime

  SetOutPath $INSTDIR

  # Copy the main binary and all Wails assets
  !insertmacro wails.files

  # Desktop shortcut
  CreateShortcut "$DESKTOP\Ambients.lnk" "$INSTDIR\Ambients.exe" \
    "" "$INSTDIR\Ambients.exe" 0 \
    SW_SHOWNORMAL "" "Ambients ambient overlay"

  # Start Menu shortcut (under its own folder so it's easy to find)
  CreateDirectory "$SMPROGRAMS\Ambients"
  CreateShortcut  "$SMPROGRAMS\Ambients\Ambients.lnk" "$INSTDIR\Ambients.exe" \
    "" "$INSTDIR\Ambients.exe" 0 \
    SW_SHOWNORMAL "" "Ambients ambient overlay"
  CreateShortcut  "$SMPROGRAMS\Ambients\Uninstall Ambients.lnk" \
    "$INSTDIR\uninstall.exe"

  # Write uninstaller and Add/Remove Programs entry
  !insertmacro wails.writeUninstaller

SectionEnd

# ─── Uninstall section ────────────────────────────────────────────────────────
Section "Uninstall"

  !insertmacro wails.setShellContext

  # Remove WebView2 user data directory
  RMDir /r "$APPDATA\Ambients"

  # Remove install directory
  RMDir /r $INSTDIR

  # Remove shortcuts
  Delete "$DESKTOP\Ambients.lnk"
  Delete "$SMPROGRAMS\Ambients\Ambients.lnk"
  Delete "$SMPROGRAMS\Ambients\Uninstall Ambients.lnk"
  RMDir  "$SMPROGRAMS\Ambients"

  # Remove Add/Remove Programs entry
  !insertmacro wails.deleteUninstaller

  # NOTE: We intentionally do NOT delete %APPDATA%\ambientspace\
  # That folder contains the user's settings, skins, and background choices.
  # Users can delete it manually if they want a clean uninstall.

SectionEnd
