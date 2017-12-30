defmodule TGBot.Pictures.Concatenators.ImageMagick do
  @resources_dir :code.priv_dir(:tg_bot)
  @tmp_dir Path.join(@resources_dir, "tmp_files")
  @glue_image Path.join(@resources_dir, "glue_gap.jpg")
  @version "v1"
  require Logger
  alias TGBot.Pictures.Concatenator
  @behaviour Concatenator

  use Utils.Meter

  @spec version :: String.t
  def version, do: @version

  @spec concatenate(String.t, String.t) :: String.t
  def concatenate(left_picture_url, right_picture_url) do
    Logger.info("Concatenate #{left_picture_url} and #{right_picture_url}")
    measure metric_name: "concatenate_photos" do
      current_dir = new_tmp_dir_path()
      try do
        Logger.info("Create tmp dir for the concatenation: #{current_dir}")
        File.mkdir!(current_dir)
        [left_picture_path, right_picture_path] = measure metric_name: "download_photos" do
          Utils.parallelize_tasks(
            [
              fn -> download_file(left_picture_url, current_dir) end,
              fn -> download_file(right_picture_url, current_dir) end
            ]
          )
        end
        {left_picture_path, right_picture_path} = ensure_same_height(
          left_picture_path,
          right_picture_path,
          current_dir
        )
        result_file_path = new_tmp_file_path()
        execute_cmd(
          "convert",
          ["+append", left_picture_path, @glue_image, right_picture_path, result_file_path]
        )
        result_file_path
      after
        Logger.info("Delete current tmp dir: #{current_dir}")
        File.rm_rf!(current_dir)
      end
    end
  end

  @spec ensure_same_height(String.t, String.t, String.t) :: {String.t, String.t}
  defp ensure_same_height(left_picture, right_picture, current_dir) do
    [left_picture_height, right_picture_height] = Utils.parallelize_tasks(
      [
        fn -> get_height(left_picture) end,
        fn -> get_height(right_picture) end
      ]
    )
    cond do
      left_picture_height == right_picture_height ->
        Logger.info("Photos height is the same, no need to crop")
        {left_picture, right_picture}
      left_picture_height < right_picture_height ->
        {left_picture, crop(right_picture, right_picture_height, left_picture_height, current_dir)}
      left_picture_height > right_picture_height ->
        {crop(left_picture, left_picture_height, right_picture_height, current_dir), right_picture}
    end
  end

  @spec crop(String.t, integer, integer, String.t) :: String.t
  defp crop(picture_uri, picture_height, result_height, current_dir) do
    crop_height = div(picture_height - result_height, 2)
    out_path = new_tmp_file_path(dir: current_dir)
    execute_cmd("convert", [picture_uri, "-crop", "x#{result_height}+0+#{crop_height}", out_path])
    out_path
  end

  @spec new_tmp_file_path(Keyword.t) :: String.t
  defp new_tmp_file_path(opts \\ []) do
    dir = Keyword.get(opts, :dir, @tmp_dir)
    Path.join(dir, UUID.uuid4() <> ".jpg")
  end

  @spec new_tmp_dir_path :: String.t
  defp new_tmp_dir_path do
    Path.join(@tmp_dir, UUID.uuid4())
  end

  @spec get_height(String.t) :: integer
  defp get_height(picture_uri) do
    result_data = execute_cmd("magick", ["identify", "-ping", "-format", "%h", picture_uri])
    case Integer.parse(result_data) do
      {height, ""} -> height
      _ -> raise "Get height command returned non-int result: #{result_data}"
    end
  end

  @spec download_file(String.t, String.t) :: String.t
  defp download_file(url, current_dir) do
    {body, _} = Utils.download_file(url, http_pool: :pictures_downloading)
    new_file_path = new_tmp_file_path(dir: current_dir)
    File.write!(new_file_path, body)
    new_file_path
  end

  @spec execute_cmd(String.t, [String.t]) :: String.t
  defp execute_cmd(cmd_name, args) do
    case System.cmd(cmd_name, args, stderr_to_stdout: true) do
      {data, 0} ->
        data
      {error_msg, error_code} ->
        raise "#{cmd_name} command execution failed, code: #{error_code}, msg: #{error_msg}"
    end
  end
end

