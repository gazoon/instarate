defmodule Deploy do
  @tmp_dir "tmp_deps"
  @local_dep_indicator ~s(, path: ")
  require Logger

  def start(_working_dir = nil), do: raise "Empty working dir"
  def start(working_dir) do
    File.cd!(working_dir)
    original_lines = "mix.exs"
                     |> File.stream!()
                     |> Enum.to_list
    create_tmp_dir()
    create_mix_file_backup()
    original_deps = get_dep_paths(original_lines)
    copy_deps(original_deps)
    patch_mix_file(original_lines)
    |> save_mix_file()
  end

  def get_dep_paths(lines) do
    lines
    |> Enum.filter(fn line -> String.contains?(line, @local_dep_indicator) end)
    |> Enum.map(fn line -> get_dependency_path(line) end)
  end

  def create_mix_file_backup() do
    File.cp!("mix.exs", "mix.exs.backup")
  end

  def copy_deps(dep_paths) do
    dep_paths
    |> Enum.each(
         fn dep_path ->
           Logger.info("copy dependecy #{dep_path}")
           dep_name = Path.basename(dep_path)
           File.cp_r!(dep_path, Path.join(@tmp_dir, dep_name))
         end
       )
  end


  def save_mix_file(lines) do
    File.write("mix.exs", Enum.join(lines, ""))
  end

  def patch_mix_file(lines) do
    lines
    |> Enum.map(
         fn line ->
           if String.contains?(line, @local_dep_indicator) do
             dep_path = get_dependency_path(line)
             dep_name = Path.basename(dep_path)
             new_dep_path = Path.join(@tmp_dir, dep_name)
             String.replace(
               line,
               @local_dep_indicator <> dep_path,
               @local_dep_indicator <> new_dep_path
             )
           else
             line
           end
         end
       )
  end

  def get_dependency_path(line) do
    Enum.at(String.split(line, ~s(")), 1)
  end

  def create_tmp_dir do
    File.mkdir!(@tmp_dir)
  end
end

Deploy.start(List.first(System.argv))
