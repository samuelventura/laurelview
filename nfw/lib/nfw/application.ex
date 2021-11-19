defmodule Nfw.Application do
  # See https://hexdocs.pm/elixir/Application.html
  # for more information on OTP Applications
  @moduledoc false

  use Application

  @impl true
  def start(_type, _args) do
    # See https://hexdocs.pm/elixir/Supervisor.html
    # for other strategies and supported options
    opts = [strategy: :one_for_one, name: Nfw.Supervisor]

    children =
      [
        # Children for all targets
        # Starts a worker by calling: Nfw.Worker.start_link(arg)
        # {Nfw.Worker, arg},
      ] ++ children(target())

    Supervisor.start_link(children, opts)
  end

  # List all child processes to be supervised
  def children(:host) do
    [
      # Children that only run on the host
      # Starts a worker by calling: Nfw.Worker.start_link(arg)
      # {Nfw.Worker, arg},
    ]
  end

  def children(_target) do
    bin = Application.fetch_env!(:nfw, :bin)
    [
      # Children for all targets except host
      # Starts a worker by calling: Nfw.Worker.start_link(arg)
      # {Nfw.Worker, arg},
      %{
        id: :lvdpm,
        start: {NervesBackdoor.Daemon, :start_link, [Path.join(bin, "lvdpm")]}
      },
      %{
        id: :lvnbe,
        start: {NervesBackdoor.Daemon, :start_link, [Path.join(bin, "lvnbe")]}
      },
      %{
        id: :lvnup,
        start: {NervesBackdoor.Daemon, :start_link, [Path.join(bin, "lvnup")]}
      },
    ]
  end

  def target() do
    Application.get_env(:nfw, :target)
  end
end
