if exist lvsss.exe (
    netsh advfirewall firewall delete rule name="LaurelViewSetup.Web"
    netsh advfirewall firewall delete rule name="LaurelViewSetup.Backend"
    net stop LaurelViewSetup
    lvsss -service uninstall
)
