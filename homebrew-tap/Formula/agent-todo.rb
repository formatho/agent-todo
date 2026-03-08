class AgentTodo < Formula
  desc "CLI tool for managing Agent Todo Platform"
  homepage "https://github.com/formatho/agent-todo"
  version "0.1.0"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/formatho/agent-todo/releases/download/v#{version}/agent-todo-darwin-amd64"
      sha256 "PLACEHOLDER"
    end
    on_arm do
      url "https://github.com/formatho/agent-todo/releases/download/v#{version}/agent-todo-darwin-arm64"
      sha256 "PLACEHOLDER"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/formatho/agent-todo/releases/download/v#{version}/agent-todo-linux-amd64"
      sha256 "PLACEHOLDER"
    end
    on_arm do
      url "https://github.com/formatho/agent-todo/releases/download/v#{version}/agent-todo-linux-arm64"
      sha256 "PLACEHOLDER"
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
