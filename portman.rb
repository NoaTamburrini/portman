class Portman < Formula
  desc "Developer-friendly CLI tool for managing ports and processes"
  homepage "https://github.com/NoaTamburrini/portman"
  url "https://github.com/NoaTamburrini/portman/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "" # Will be filled after first release
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w")
  end

  test do
    assert_match "Portman - Port Management CLI Tool", shell_output("#{bin}/portman help")
  end
end
