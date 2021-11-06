defmodule Nss.Daemon do
  use GenServer

  @spawn_opts [:stderr_to_stdout, :binary, :stream, {:line, 255}]

  #Nss.Setup.setup
  #{:ok, pid} = Nss.Daemon.start_link "/tmp/lvbin/lvdpm"
  #:ok = GenServer.stop pid
  def start_link(path) do
    GenServer.start_link(__MODULE__, path)
  end

  defp run(path) do
    env = {:env, Nss.Setup.env path}
    opts = [env | @spawn_opts]
    IO.inspect {"run", path, env}
    port = Port.open({:spawn_executable, path}, opts)
    _ref = Port.monitor(port)
    port
  end

  @impl true
  def init(path) do
    IO.inspect {"init", path}
    port = run path
    IO.inspect {"state", {path, port}}
    {:ok, {path, port}}
  end

  @impl true
  def terminate(reason, state={_path, port}) do
    IO.inspect {"terminate", reason, state}
    :true = Port.close port
  end

  @impl true
  def handle_info({_port, {:data, data}}, state) do
    IO.inspect {state, data}
    {:noreply, state}
  end

  @impl true
  def handle_info({:DOWN, _ref, :port, port, reason}, state={path, _}) do
    IO.inspect {state, "down", port, reason}
    port = run path
    IO.inspect {"state", {path, port}}
    {:noreply, {path, port}}
  end
end
