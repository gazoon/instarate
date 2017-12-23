defmodule TGBot.Application do

  use Application

  require Logger

  def start(_type, _args) do
    Logger.info "Started application tg_bot"
    TGBot.Supervisor.start_link([])
    #    stuff()
  end

  #  def stuff do
  #    Process.flag(:trap_exit, true)
  #    try do
  #      try do
  #        IO.inspect("start task")
  #        t = Task.async(
  #          #    t = Task.Supervisor.async_nolink(
  #          #      :message_workers_supervisor,
  #          fn ->
  #            Process.sleep(1000)
  #            raise "ffff"
  #          end
  #        )
  #        IO.inspect("start awaiting")
  #        Task.await(t)
  #      catch
  #        :exit, {{error, stack}, from} ->
  #          IO.inspect("catch error")
  #          reraise error, stack
  #      end
  #    rescue
  #      e -> IO.inspect("rescue error")
  #           IO.inspect(e)
  #    after
  #      IO.inspect("after")
  #    end
  #    IO.inspect("success ending")
  #  end

end
