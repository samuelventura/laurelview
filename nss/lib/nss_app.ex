defmodule Nss.App do
  use Application

  #cat _build/dev/lib/nss/ebin/nss.app
  #Application.start :nss
  #Application.stop :nss
  @impl true
  def start(_type, _args) do
    Nss.Super.start_link(name: Nss.Super)
  end
end
