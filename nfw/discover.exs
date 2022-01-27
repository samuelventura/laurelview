# mix run discover.exs --no-start

loop = fn socket, loop ->
  case :gen_udp.recv(socket, 1024, 1000) do
    {:ok, msg} ->
      IO.inspect(msg)
      loop.(socket, loop)

    {:error, msg} ->
      IO.inspect(msg)
  end
end

cmd = Jason.encode!(%{action: "id", name: "lvbox"})

{:ok, socket} =
  :gen_udp.open(0, active: false, broadcast: true, multicast_ttl: 4, multicast_loop: false)

:ok = :gen_udp.send(socket, {255, 255, 255, 255}, 31680, cmd)
loop.(socket, loop)
:ok = :gen_udp.close(socket)
