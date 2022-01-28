rem netsh advfirewall firewall show rule name="LaurelView.Backend"
netsh advfirewall firewall add rule name="LaurelViewSetup.Backend" dir=in action=allow program="%~dp0lvsbe.exe" enable=yes
rem netsh advfirewall firewall show rule name="LaurelViewSetup.Web"
netsh advfirewall firewall add rule name="LaurelViewSetup.Web" dir=in action=allow protocol=TCP localport=31603
lvsss -service install
net start LaurelViewSetup
