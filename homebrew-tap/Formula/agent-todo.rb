class AgentTodo < Formula
  desc "CLI tool for managing Agent Todo Platform"
  homepage "https://github.com/formatho/agent-todo"
  version "0.1.0"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/formatho/agent-todo/releases/download/v#{version}/agent-todo-darwin-amd64"
      sha256 "b986a35db9bdd8d4a2fb4822758fd976a0bbf242841628005eb7575c82888ec4"
    end
    on_arm do
      url "https://github.com/formatho/agent-todo/releases/download/v#{version}/agent-todo-darwin-arm64"
      sha256 "6595bb9462e0dc34319b772dfdea56031157bfdb098c953168f1f78928033817"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/formatho/agent-todo/releases/download/v#{version}/agent-todo-linux-amd64"
      sha256 "a0d108503e96882265a84129a423050ed372a225453e3b8209855f2f2c8735c6"
    end
    on_arm do
      url "https://github.com/formatho/agent-todo/releases/download/v#{version}/agent-todo-linux-arm64"
      sha256 "69ed8aa214ed62dd79f619acf97d5ddbf85b37f617df316d9969ae381d1ac472"
    end
  end

  def install
    binary_name = "agent-todo-#{OS.kernel_name.downcase}-#{Hardware::CPU.intel? ? 'amd64' : 'arm64'}"
    bin.install binary_name => "agent-todo"
  end

  test do
    system "#{bin}/agent-todo", "version"
  end
end
