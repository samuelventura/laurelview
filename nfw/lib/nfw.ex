defmodule Nfw do
  def setup do
    usb = Application.fetch_env!(:nfw, :usb)
    System.cmd "mkdir", ["-p", "/tmp/usb"]
    if File.exists?(usb) do
      IO.inspect "Found #{usb}"
      System.cmd "mount", [usb, "/tmp/usb"]
      path = "/tmp/usb/lvnet.txt"
      if File.exists?(path) do
        IO.inspect "Found #{path}"
        case File.read(path) do
          {:ok, data} ->
            IO.inspect "Net config #{data}"
            case parse(data) do
              {:ok, term} -> VintageNet.configure("eth0", term)
              _ -> IO.inspect "Invalid #{data}"
            end
        end
      end
      System.cmd "umount", ["/tmp/usb"]
      end
      #VintageNet.info
      #VintageNet.configure("eth0", %{type: VintageNetEthernet, ipv4: %{method: :dhcp}})
      #VintageNet.configure("eth0", %{type: VintageNetEthernet, ipv4: %{ method: :static, address: "10.77.4.165", prefix_length: 8, gateway: "10.77.0.1", name_servers: ["10.77.0.1"]}})
  end

  def parse(str) when is_binary(str) do
    case str |> Code.string_to_quoted do
      {:ok, terms} -> {:ok, _parse(terms)}
      {:error, _}  -> {:invalid_terms}
    end
  end

  # atomic terms
  defp _parse(term) when is_atom(term), do: term
  defp _parse(term) when is_integer(term), do: term
  defp _parse(term) when is_float(term), do: term
  defp _parse(term) when is_binary(term), do: term

  defp _parse([]), do: []
  defp _parse([h|t]), do: [_parse(h) | _parse(t)]

  defp _parse({a, b}), do: {_parse(a), _parse(b)}
  defp _parse({:"{}", _place, terms}) do
    terms
    |> Enum.map(&_parse/1)
    |> List.to_tuple
  end

  defp _parse({:"%{}", _place, terms}) do
    for {k, v} <- terms, into: %{}, do: {_parse(k), _parse(v)}
  end

  defp _parse({_term_type, _place, terms}), do: terms # to ignore functions and operators

end
