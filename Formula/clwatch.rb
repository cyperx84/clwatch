class Clwatch < Formula
  desc "Track coding tool updates from changelogs.info"
  homepage "https://github.com/cyperx84/clwatch"
  version "1.0.0"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/cyperx84/clwatch/releases/download/v#{version}/clwatch-#{version}-darwin-arm64.tar.gz"
      sha256 "PLACEHOLDER_DARWIN_ARM64"
    end

    on_intel do
      url "https://github.com/cyperx84/clwatch/releases/download/v#{version}/clwatch-#{version}-darwin-amd64.tar.gz"
      sha256 "PLACEHOLDER_DARWIN_AMD64"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/cyperx84/clwatch/releases/download/v#{version}/clwatch-#{version}-linux-arm64.tar.gz"
      sha256 "PLACEHOLDER_LINUX_ARM64"
    end

    on_intel do
      url "https://github.com/cyperx84/clwatch/releases/download/v#{version}/clwatch-#{version}-linux-amd64.tar.gz"
      sha256 "PLACEHOLDER_LINUX_AMD64"
    end
  end

  def install
    bin.install "clwatch"
  end

  test do
    assert_match "clwatch #{version}", shell_output("#{bin}/clwatch version")
  end
end
