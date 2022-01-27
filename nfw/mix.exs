defmodule Nfw.MixProject do
  use Mix.Project

  @app :nfw
  @version "0.2.0"
  @all_targets [:bbb, :bbb_emmc, :rpi4]

  def project do
    [
      app: @app,
      version: @version,
      elixir: "~> 1.9",
      archives: [nerves_bootstrap: "~> 1.10"],
      start_permanent: Mix.env() == :prod,
      build_embedded: true,
      deps: deps(),
      releases: [{@app, release()}],
      preferred_cli_target: [run: :host, test: :host]
    ]
  end

  # Run "mix help compile.app" to learn about applications.
  def application do
    [
      mod: {Nfw.Application, []},
      extra_applications: [
        :logger,
        :runtime_tools,
        :plug,
        :jason,
        :plug_cowboy,
        :mac_address,
        :circuits_gpio
      ],
      # usb is partition to mount
      env: [
        usb: System.get_env("NFW_USB") || "/dev/sda1",
        bin: System.get_env("NFW_BIN") || "/lvbin",
        port: 31680,
        name: "lvbox",
        home: "/data/nfw",
        version: @version,
        ifname: "eth0",
        io_push: 47,
        io_red: 66,
        io_green: 45,
        io_blue: 69,
        blink_ms: 200,
        blink_color: :blue,
        reset_color: :red,
        restart_color: :green,
        reset_ms: 3000
      ]
    ]
  end

  # Run "mix help deps" to learn about dependencies.
  defp deps do
    [
      # Dependencies for all targets
      {:nerves, "~> 1.7.4", runtime: false},
      {:shoehorn, "~> 0.7.0"},
      {:ring_logger, "~> 0.8.1"},
      {:toolshed, "~> 0.2.13"},

      # Dependencies for all targets except :host
      {:nerves_runtime, "~> 0.11.3", targets: @all_targets},
      {:nerves_pack, "~> 0.6.0", targets: @all_targets},

      # Dependencies for specific targets
      # NOTE: It's generally low risk and recommended to follow minor version
      # bumps to Nerves systems. Since these include Linux kernel and Erlang
      # version updates, please review their release notes in case
      # changes to your application are needed.
      {:nerves_system_rpi4, "~> 1.17", runtime: false, targets: :rpi4},
      {:nerves_system_bbb, "~> 2.12", runtime: false, targets: :bbb},
      {:nerves_system_bbb_emmc, "~> 0.0.1", runtime: false, targets: :bbb_emmc},
      {:plug, "~> 1.7"},
      {:jason, "~> 1.2"},
      {:plug_cowboy, "~> 2.5"},
      {:mac_address, "~> 0.0.1"},
      {:circuits_gpio, "~> 0.4"},
      {:cors_plug, "~> 2.0"},
      {:ex_doc, ">= 0.0.0", only: :dev, runtime: false}
    ]
  end

  def release do
    [
      overwrite: true,
      # Erlang distribution is not started automatically.
      # See https://hexdocs.pm/nerves_pack/readme.html#erlang-distribution
      cookie: "#{@app}_cookie",
      include_erts: &Nerves.Release.erts/0,
      steps: [&Nerves.Release.init/1, :assemble],
      strip_beams: Mix.env() == :prod or [keep: ["Docs"]]
    ]
  end
end
