rem netsh advfirewall firewall show rule name="LaurelView.Web"
netsh advfirewall firewall add rule name="LaurelView.Web" dir=in action=allow program="%~dp0lvnrt.exe" enable=yes
rem netsh advfirewall firewall show rule name="LaurelView.Runtime"
netsh advfirewall firewall add rule name="LaurelView.Runtime" dir=in action=allow program="%~dp0lvnbe.exe" enable=yes
rem netsh advfirewall firewall show rule name="LaurelView.Web"
netsh advfirewall firewall add rule name="LaurelView.Web" dir=in action=allow protocol=TCP localport=31601
rem netsh advfirewall firewall show rule name="LaurelView.Runtime"
netsh advfirewall firewall add rule name="LaurelView.Runtime" dir=in action=allow protocol=TCP localport=31602
lvnss -service install
net start LaurelView
