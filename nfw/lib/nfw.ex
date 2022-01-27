defmodule Nfw do
  def env_home() do
    Application.get_env(:nfw, :home)
  end

  def env_name() do
    Application.get_env(:nfw, :name)
  end

  def env_hostname() do
    Application.get_env(:nfw, :hostname, get_hostname())
  end

  def env_version() do
    Application.get_env(:nfw, :version)
  end

  def env_port() do
    Application.get_env(:nfw, :port)
  end

  def env_ifname() do
    Application.get_env(:nfw, :ifname)
  end

  def env_red() do
    Application.get_env(:nfw, :io_red)
  end

  def env_green() do
    Application.get_env(:nfw, :io_green)
  end

  def env_blue() do
    Application.get_env(:nfw, :io_blue)
  end

  def env_push() do
    Application.get_env(:nfw, :io_push)
  end

  def env_blink_ms() do
    Application.get_env(:nfw, :blink_ms)
  end

  def env_blink_color() do
    Application.get_env(:nfw, :blink_color)
  end

  def env_reset_ms() do
    Application.get_env(:nfw, :reset_ms)
  end

  def env_reset_color() do
    Application.get_env(:nfw, :reset_color)
  end

  def env_restart_color() do
    Application.get_env(:nfw, :restart_color)
  end

  def get_mac(ifname \\ :env) do
    ifname =
      case ifname do
        :env -> env_ifname()
        ifname -> ifname
      end

    case MACAddress.mac_address(ifname) do
      {:ok, mac} -> mac |> MACAddress.to_hex(case: :upper)
      _ -> "00:00:00:00:00:00"
    end
    |> String.replace(":", "")
  end

  def get_pass(type \\ :current) do
    case type do
      :default ->
        get_mac()

      :current ->
        path = pass_path()

        case File.read(path) do
          {:ok, data} -> String.trim(data)
          _ -> get_mac()
        end
    end
  end

  def reset_pass() do
    File.rm(pass_path())
  end

  def disable_pass() do
    set_pass("")
  end

  def set_pass(password) do
    File.write(pass_path(), password)
  end

  def pass_path() do
    Path.join(env_home(), "password.txt")
  end

  def get_hostname() do
    {:ok, hostname} = :inet.gethostname()
    hostname |> to_string
  end

  def io_blink(port, ms \\ :env) do
    io =
      case port do
        :red -> Nfw.env_red()
        :green -> Nfw.env_green()
        :blue -> Nfw.env_blue()
        port -> port
      end

    ms =
      case ms do
        :env -> Nfw.env_blink_ms()
        ms -> ms
      end

    {:ok, gpio} = Nfw.Gpio.output(io)
    :ok = Nfw.Gpio.write(gpio, 1)
    :timer.sleep(ms)
    :ok = Nfw.Gpio.write(gpio, 0)
    :timer.sleep(ms)
    :ok = Nfw.Gpio.write(gpio, 1)
    :timer.sleep(ms)
    :ok = Nfw.Gpio.write(gpio, 0)
    :timer.sleep(ms)
    :ok = Nfw.Gpio.write(gpio, 1)
    :timer.sleep(ms)
    :ok = Nfw.Gpio.write(gpio, 0)
    :ok = Nfw.Gpio.close(gpio)
  end

  def io_input(port \\ :env) do
    io =
      case port do
        :env -> Nfw.env_push()
        :push -> Nfw.env_push()
        :red -> Nfw.env_red()
        :green -> Nfw.env_green()
        :blue -> Nfw.env_blue()
        port -> port
      end

    {:ok, gpio} = Nfw.Gpio.input(io)
    value = Nfw.Gpio.read(gpio)
    :ok = Nfw.Gpio.close(gpio)
    value
  end

  def io_output(port, value) do
    io =
      case port do
        :red -> Nfw.env_red()
        :green -> Nfw.env_green()
        :blue -> Nfw.env_blue()
        port -> port
      end

    {:ok, gpio} = Nfw.Gpio.output(io)
    Nfw.Gpio.write(gpio, value)
    :ok = Nfw.Gpio.close(gpio)
    value
  end

  def reset_network() do
    Nfw.Vintage.configure(
      Nfw.env_ifname(),
      %{"method" => "dhcp"}
    )
  end
end
