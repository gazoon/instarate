defmodule TGBot.Pictures.Concatenator do
  @callback concatenate(left_picture_url :: String.t, right_picture_url :: String.t) :: String.t
  @callback version :: String.t
end
