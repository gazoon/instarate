defmodule TGBot.Pictures do

  @resources_dir :code.priv_dir(:tg_bot)
  @tmp_dir Path.join(@resources_dir, "tmp_files")
  File.mkdir(@tmp_dir)
  @glue_image Path.join(@resources_dir, "glue_gap.jpg")
  @version "v1"
  require Logger

  def version, do: @version

  @spec concatenate(String.t, String.t) :: String.t
  def concatenate(left_picture, right_picture) do
    Logger.info("Concatenate #{left_picture} and #{right_picture}")
    {left_picture, right_picture} = ensure_same_height(left_picture, right_picture)
    try do
      result_file_path = new_tmp_file_path()
      execute_cmd(
        "convert",
        ["+append", left_picture, @glue_image, right_picture, result_file_path]
      )
      result_file_path
    after
      clean_tmp_files([left_picture, right_picture])
    end
  end

  @spec clean_tmp_files([String.t]) :: any
  defp clean_tmp_files(files) do
    Enum.each(files, &clean_tmp_file/1)
  end

  @spec clean_tmp_file(String.t) :: any
  defp clean_tmp_file(file_path) do
    if String.starts_with?(file_path, @tmp_dir)  do
      case File.rm(file_path) do
        {:error, details} -> Logger.warn("Can't delete tmp file #{file_path}: #{details}")
        _ -> nil
      end
    end
  end

  @spec ensure_same_height(String.t, String.t) :: {String.t, String.t}
  defp ensure_same_height(left_picture, right_picture) do
    task = Task.async(fn -> get_height(left_picture) end)
    right_picture_height = get_height(right_picture)
    left_picture_height = Task.await(task)
    cond do
      left_picture_height == right_picture_height ->
        {left_picture, right_picture}
      left_picture_height < right_picture_height ->
        {left_picture, crop(right_picture, right_picture_height, left_picture_height)}
      left_picture_height > right_picture_height ->
        {crop(left_picture, left_picture_height, right_picture_height), right_picture}
    end
  end

  @spec crop(String.t, integer, integer) :: String.t
  defp crop(picture_uri, picture_height, result_height) do
    crop_height = div(picture_height - result_height, 2)
    out_path = new_tmp_file_path()
    execute_cmd("convert", [picture_uri, "-crop", "x#{result_height}+0+#{crop_height}", out_path])
    out_path
  end

  @spec new_tmp_file_path :: String.t
  defp new_tmp_file_path do
    Path.join(@tmp_dir, UUID.uuid4() <> ".jpg")
  end

  @spec get_height(String.t) :: integer
  defp get_height(picture_uri) do
    result_data = execute_cmd("magick", ["identify", "-ping", "-format", "%h", picture_uri])
    case Integer.parse(result_data) do
      {height, ""} -> height
      _ -> raise "Get height command returned non-int result: #{result_data}"
    end
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
