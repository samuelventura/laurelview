if exist lvnss.exe (
    netsh advfirewall firewall delete rule name="LaurelView.Web"
    netsh advfirewall firewall delete rule name="LaurelView.Backend"
    net stop LaurelView
    lvnss -service uninstall
)
