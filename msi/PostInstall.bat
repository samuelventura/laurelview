rem netsh advfirewall firewall show rule name="LaurelView.Backend"
netsh advfirewall firewall add rule name="LaurelView.Backend" dir=in action=allow program="%~dp0lvnbe.exe" enable=yes
rem netsh advfirewall firewall show rule name="LaurelView.Web"
netsh advfirewall firewall add rule name="LaurelView.Web" dir=in action=allow protocol=TCP localport=31601
lvnss -service install
net start LaurelView
