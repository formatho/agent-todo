class AgentTodo < Formula
  desc "CLI tool for managing Agent Todo Platform"
  homepage "https://github.com/formatho/agent-todo"
  version "0.1.0"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/formatho/agent-todo/releases/download/v#{version}/agent-todo-darwin-amd64"
      sha256 "08b6f243a7cc585aa88abdea9be61de403cb24e22adfd4382dd333e6ea0f7098"
    end
    on_arm do
      url "https://github.com/formatho/agent-todo/releases/download/v#{version}/agent-todo-darwin-arm64"
      sha256 "ec2cf743abd39ed8676caa1782aeb940982c3f1d8e99cc357c91e677ba3b5866"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/formatho/agent-todo/releases/download/v#{version}/agent-todo-linux-amd64"
      sha256 "764d6adb68e8b457ab792da8b4f7d0aebc403885892fe46da4f585bdb1d11da8"
    end
    on_arm do
      url "https://github.com/formatho/agent-todo/releases/download/v#{version}/agent-todo-linux-arm64"
      sha256 "6a201de6864de21898fa6b3e4117a0b64489f61f1844e675953d94a9993b692b"
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
