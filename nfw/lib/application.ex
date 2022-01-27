defmodule Nfw.Application do
  use Application

  @impl true
  def start(_type, _args) do
    bin = Application.fetch_env!(:nfw, :bin)
    port = Nfw.env_port()
    home = Nfw.env_home()
    File.mkdir_p(home)
    off(Nfw.env_red())
    off(Nfw.env_green())
    off(Nfw.env_blue())

    children = [
      {Nfw.Reset, []},
      {Nfw.Gpio, []},
      {Nfw.Vintage, []},
      {Nfw.Discovery, []},
      Plug.Cowboy.child_spec(
        scheme: :http,
        plug: Nfw.Endpoint,
        options: [port: port]
      ),
      %{
        id: :lvdpm,
        start: {Nfw.Daemon, :start_link, [Path.join(bin, "lvdpm")]}
      },
      %{
        id: :lvnbe,
        start: {Nfw.Daemon, :start_link, [Path.join(bin, "lvnbe")]}
      },
      %{
        id: :lvnup,
        start: {Nfw.Daemon, :start_link, [Path.join(bin, "lvnup")]}
      }
    ]

    opts = [strategy: :one_for_one, name: Nfw.Supervisor]
    Supervisor.start_link(children, opts)
  end

  defp off(port) do
    {:ok, gpio} = Nfw.Gpio.io_output(port)
    :ok = Nfw.Gpio.io_write(gpio, 0)
    :ok = Nfw.Gpio.io_close(gpio)
  end
end
