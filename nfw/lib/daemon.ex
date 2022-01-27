defmodule Nfw.Daemon do
  use GenServer

  @spawn_opts [:stderr_to_stdout, :binary, :stream, {:line, 255}]

  # {:ok, pid} = Nfw.Daemon.start_link "/bin/daemon"
  # :ok = GenServer.stop pid
  def start_link(path) do
    GenServer.start_link(__MODULE__, path)
  end

  @impl true
  def init(path) do
    port = run(path)
    {:ok, {path, port}}
  end

  @impl true
  def terminate(reason, _state = {path, port}) do
    IO.inspect({path, :terminate, reason})
    true = Port.close(port)
  end

  @impl true
  def handle_info({_port, {:data, data}}, state = {path, _}) do
    IO.inspect({path, :data, data})
    {:noreply, state}
  end

  @impl true
  def handle_info({:DOWN, _ref, :port, _port, reason}, _state = {path, _}) do
    IO.inspect({path, :down, reason})
    port = run(path)
    {:noreply, {path, port}}
  end

  defp run(path) do
    env = {:env, env(path)}
    opts = [env | @spawn_opts]
    IO.inspect({path, :run, env})
    port = Port.open({:spawn_executable, path}, opts)
    _ref = Port.monitor(port)
    port
  end

  defp env(path) do
    case File.read(path <> ".env") do
      {:ok, data} ->
        data
        |> String.split("\n")
        |> Enum.filter(&filter/1)
        |> Enum.map(&tuple/1)

      _ ->
        []
    end
  end

  defp filter(line) do
    String.contains?(line, "=")
  end

  defp tuple(line) do
    trimmed = String.trim(line)
    [n, v] = String.split(trimmed, "=", parts: 2)
    {String.to_charlist(n), String.to_charlist(v)}
  end
end
