defmodule Instagram.Media do

  alias Instagram.Media
  @type t :: %Media{owner: String.t, url: String.t, is_photo: boolean}

  defstruct owner: nil, url: nil, is_photo: true

end
