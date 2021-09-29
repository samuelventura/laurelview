; Script generated by the Inno Setup Script Wizard.
; SEE THE DOCUMENTATION FOR DETAILS ON CREATING INNO SETUP SCRIPT FILES!

#define MyAppId "LaurelView" 
#define MyAppName "Laurel View"
#define MyAppVersion "0.0.4"
#define MyAppPublisher "Samuel Ventura"
#define MyAppURL "https://github.com/samuelventura/laurelview"

[Setup]
; NOTE: The value of AppId uniquely identifies this application.
; Do not use the same AppId value in installers for other applications.
; (To generate a new GUID, click Tools | Generate GUID inside the IDE.)
AppId={{C7627C4A-EC39-41E5-9712-755714B1C393}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppVerName={#MyAppName} {#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}
DefaultDirName={pf}\{#MyAppPublisher}\{#MyAppName}
DefaultGroupName={#MyAppName}
OutputBaseFilename={#MyAppId}-{#MyAppVersion}
SetupIconFile=icon.ico
Compression=lzma
SolidCompression=yes
UninstallDisplayIcon={uninstallexe}
ChangesAssociations = yes
OutputDir=build

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Files]
Source: "c:\Users\samuel\go\bin\lv*.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "msi\*.bat"; DestDir: "{app}"; Flags: ignoreversion
Source: "icon.ico"; DestDir: "{app}"; Flags: ignoreversion
; NOTE: Don't use "Flags: ignoreversion" on any shared system files
                                                
[Icons]
Name: "{group}\{#MyAppName}"; Filename: "http://127.0.0.1:31601"; IconFilename: "{app}\icon.ico"
Name: "{commondesktop}\{#MyAppName}"; Filename: "http://127.0.0.1:31601"; IconFilename: "{app}\icon.ico"
;Version in icon leaves previous link when upgrading

[Run]
Filename: "{app}\PostInstall.bat";

[UninstallRun]
Filename: "{app}\PreUninstall.bat";

[Code]
function PrepareToInstall(var NeedsRestart: Boolean): String;
var
  ResultCode: integer;
begin
  Exec(ExpandConstant('{app}\PreUninstall.bat'), '', '', SW_SHOW, ewWaitUntilTerminated, ResultCode)
end;
