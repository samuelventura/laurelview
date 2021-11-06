defmodule Nss.Super do
  use Supervisor

  #{:ok, pid} = Nss.Super.start_link
  #:ok = Supervisor.stop pid
  def start_link(opts \\ []) do
    Supervisor.start_link(__MODULE__, :ok, opts)
  end

  @impl true
  def init(:ok) do
    folder = Nss.Setup.folder

    children = [
      %{
        id: :lvdpm,
        start: {Nss.Daemon, :start_link, [Path.join(folder, "lvdpm")]}
      },
      %{
        id: :lvnbe,
        start: {Nss.Daemon, :start_link, [Path.join(folder, "lvnbe")]}
      },
      %{
        id: :lvnup,
        start: {Nss.Daemon, :start_link, [Path.join(folder, "lvnup")]}
      },
    ]

    Supervisor.init(children, strategy: :one_for_one)
  end
end
