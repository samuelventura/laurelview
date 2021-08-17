if exist lvnss.exe (
    netsh advfirewall firewall delete rule name="LaurelView.Web"
    netsh advfirewall firewall delete rule name="LaurelView.Runtime"
    net stop LaurelView
    lvnss -service uninstall
)
