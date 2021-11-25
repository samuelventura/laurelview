import Config

config :nerves, :erlinit,
  ctty: "ttyS0"

config :nerves_backdoor,
  io_red: 20,
  io_green: 20,
  io_blue: 20,
  io_push: 21
