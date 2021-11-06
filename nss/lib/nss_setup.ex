defmodule Nss.Setup do

  #Nss.Setup.setup
  def setup() do
    Application.fetch_env!(:nss, :folder)
  end

  def env(path) do
    case File.read(path <> ".env") do
      {:ok, data} ->
        data
        |> String.split("\n")
        |> Enum.filter(&filter/1)
        |> Enum.map(&tuple/1)
      _ -> []
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
