defmodule NfwTest do
  use ExUnit.Case
  doctest Nfw

  test "greets the world" do
    assert Nfw.hello() == :world
  end
end
