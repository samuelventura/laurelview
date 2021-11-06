defmodule Nss.Setup do

  #Nss.Setup.setup
  def setup() do
    folder = Application.fetch_env!(:nss, :folder)
    zip_path = Path.join :code.priv_dir(:nss), "lvbin.zip"
    {:ok, files} = :zip.unzip(String.to_charlist(zip_path), cwd: folder)
    for file <- files do
      chmod file
    end
    folder
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

  defp chmod(path) do
    if ! String.ends_with?(path, ".env") do
      :ok = File.chmod path, 0o755
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
