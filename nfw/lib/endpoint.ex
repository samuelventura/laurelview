defmodule Nfw.Endpoint do
  use Plug.Router
  use Plug.ErrorHandler

  plug(CORSPlug)
  plug(:auth)

  defp auth(conn, _opts) do
    username = Nfw.env_name()
    password = Nfw.get_pass(:current)
    # disable password check by uploading empty password file
    case String.length(password) do
      0 -> conn
      _ -> Plug.BasicAuth.basic_auth(conn, username: username, password: password)
    end
  end

  plug(:match)

  # curl http://localhost:31680/tmp/test.txt
  plug(Plug.Static, at: "/tmp", from: "/tmp")
  # curl http://localhost:31680/data/test.txt
  # curl http://nerves.local:31680/data/lvnbe.db3 --output /tmp/lvnbe.db3
  plug(Plug.Static, at: "/data", from: "/data")

  plug(Plug.Parsers, parsers: [:urlencoded, :multipart], pass: ["*/*"])
  plug(Plug.Parsers, parsers: [:json], pass: ["application/json"], json_decoder: Jason)

  # curl http://localhost:31680/ping
  # curl http://nerves.local:31680/ping
  get "/ping" do
    respond(conn, {:ok, %{ping: "pong"}})
  end

  # curl -F 'file=@/tmp/test.txt' http://localhost:31680/upload?path=/tmp/testup.txt
  # curl -F 'file=@/tmp/test.txt' http://nerves.local:31680/upload?path=/tmp/testup.txt
  post "/upload" do
    {:ok, upload} = Map.fetch(conn.params, "file")
    {:ok, path} = Map.fetch(conn.query_params, "path")
    result = File.copy(upload.path, path)
    File.rm(upload.path)
    respond(conn, result)
  end

  # curl -X DELETE http://localhost:31680/delete?path=/tmp/testup.txt
  # curl -X DELETE http://nerves.local:31680/delete?path=/tmp/testup.txt
  delete "/delete" do
    {:ok, path} = Map.fetch(conn.query_params, "path")
    result = File.rm(path)
    respond(conn, result)
  end

  # curl http://localhost:31680/pass/reset
  # curl http://nerves.local:31680/pass/reset
  get "/pass/reset" do
    result = Nfw.reset_pass()
    respond(conn, result)
  end

  # curl http://localhost:31680/pass/disable
  # curl http://nerves.local:31680/pass/disable
  get "/pass/disable" do
    result = Nfw.disable_pass()
    respond(conn, result)
  end

  # curl -X POST http://localhost:31680/pass/set --data "pass=secret64"
  # curl -X POST http://nerves.local:31680/pass/set --data "pass=secret64"
  post "/pass/set" do
    {:ok, pass64} = Map.fetch(conn.params, "pass")
    {:ok, pass} = Base.decode64(pass64)
    result = Nfw.set_pass(pass)
    respond(conn, result)
  end

  # curl http://localhost:31680/net/all
  # curl http://nerves.local:31680/net/all
  get "/net/all" do
    result = Nfw.Vintage.all_interfaces()
    respond(conn, result)
  end

  # VintageNet.info
  # VintageNet.get_configuration("eth0")
  # VintageNet.get(["interface", "eth0", "type"])
  # VintageNet.get(["interface", "eth0", "state"])
  # VintageNet.get(["interface", "eth0", "connection"])
  # curl http://localhost:31680/net/state/eth0
  # curl http://nerves.local:31680/net/state/eth0
  get "/net/state/:interface" do
    {:ok, interface} = Map.fetch(conn.path_params, "interface")
    result = Nfw.Vintage.get_configuration(interface)
    respond(conn, result)
  end

  # VintageNet.configure("eth0", %{type: VintageNetEthernet, ipv4: %{method: :dhcp}})
  # VintageNet.configure("eth0", %{type: VintageNetEthernet, ipv4: %{method: :static, address: "10.77.4.100", prefix_length: 8, gateway: "10.77.0.1", name_servers: ["10.77.0.1"]}})
  # curl -X POST http://localhost:31680/net/setup/eth0 -H "Content-Type: application/json" -d '{"method":"dhcp"}'
  # curl -X POST http://nerves.local:31680/net/setup/eth0 -H "Content-Type: application/json" -d '{"method":"dhcp"}'
  # curl -X POST http://localhost:31680/net/setup/eth0 -H "Content-Type: application/json" -d '{"method":"static", "address":"10.77.4.100", "prefix_length":8, "gateway":"10.77.0.1", "name_servers":["10.77.0.1"]}'
  # curl -X POST http://nerves.local:31680/net/setup/eth0 -H "Content-Type: application/json" -d '{"method":"static", "address":"10.77.4.100", "prefix_length":8, "gateway":"10.77.0.1", "name_servers":["10.77.0.1"]}'
  post "/net/setup/:interface" do
    {:ok, interface} = Map.fetch(conn.path_params, "interface")
    result = Nfw.Vintage.configure(interface, conn.body_params)
    respond(conn, result)
  end

  # Application.started_applications
  # Application.loaded_applications
  # Application.get_all_env :nss
  # curl -X POST http://localhost:31680/app/start/nss
  # curl -X POST http://nerves.local:31680/app/start/nss
  post "/app/start/:app" do
    {:ok, app} = Map.fetch(conn.path_params, "app")

    result =
      case app do
        "lvapp" ->
          :ok = Supervisor.terminate_child(Nfw.Supervisor, :lvnbe)
          :ok = Supervisor.terminate_child(Nfw.Supervisor, :lvnup)
          :ok

        _ ->
          {:error, "App not found #{app}"}
      end

    respond(conn, result)
  end

  # curl -X POST http://localhost:31680/app/stop/nss
  # curl -X POST http://nerves.local:31680/app/stop/nss
  post "/app/stop/:app" do
    {:ok, app} = Map.fetch(conn.path_params, "app")
    IO.inspect(app)

    result =
      case app do
        "lvapp" ->
          {:ok, _} = Supervisor.restart_child(Nfw.Supervisor, :lvnup)
          {:ok, _} = Supervisor.restart_child(Nfw.Supervisor, :lvnbe)
          :ok

        _ ->
          {:error, "App not found #{app}"}
      end

    respond(conn, result)
  end

  @impl Plug.ErrorHandler
  def handle_errors(conn, %{kind: kind, reason: reason, stack: _stack}) do
    send_resp(
      conn,
      500,
      Jason.encode!(%{result: "error", message: %{kind: kind, reason: reason}})
    )
  end

  defp app_loaded(app) do
    Enum.find(Application.loaded_applications(), fn tuple ->
      {atom, _desc, _ver} = tuple
      Atom.to_string(atom) == app
    end)
  end

  defp respond(conn, result) do
    conn = put_resp_content_type(conn, "application/json")

    case result do
      :ok ->
        send_resp(conn, 200, Jason.encode!(%{result: "ok"}))

      {:ok, message} ->
        send_resp(conn, 200, Jason.encode!(%{result: "ok", message: message}))

      {:error, message} ->
        send_resp(conn, 400, Jason.encode!(%{result: "error", message: message}))
    end
  end

  alias Jason.Encoder

  defimpl Encoder, for: Tuple do
    def encode(data, options) when is_tuple(data) do
      data
      |> Tuple.to_list()
      |> Encoder.List.encode(options)
    end
  end

  plug(:dispatch)

  match _ do
    send_resp(conn, 404, Jason.encode!(%{result: "error", message: "not found"}))
  end
end
