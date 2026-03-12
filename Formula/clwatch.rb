class Clwatch < Formula
  desc "Track coding tool updates from changelogs.info"
  homepage "https://github.com/cyperx84/clwatch"
  version "1.0.0"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/cyperx84/clwatch/releases/download/v#{version}/clwatch-#{version}-darwin-arm64.tar.gz"
      sha256 "6c0c99f2a03a2f1ed4678806bea1b9bf2d251728499b39b38219437353cf9576"
    end

    on_intel do
      url "https://github.com/cyperx84/clwatch/releases/download/v#{version}/clwatch-#{version}-darwin-amd64.tar.gz"
      sha256 "0d3022db9c922d0f2572affbe3b7e2cb26000558e02ac1c284aec1073a7737eb"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/cyperx84/clwatch/releases/download/v#{version}/clwatch-#{version}-linux-arm64.tar.gz"
      sha256 "77fd6f826143ab3eb6f1ee12f03f6c970b48baacfc0e3dcae1f6be786d111b29"
    end

    on_intel do
      url "https://github.com/cyperx84/clwatch/releases/download/v#{version}/clwatch-#{version}-linux-amd64.tar.gz"
      sha256 "b67776495065aa4a3b136e33f433ba66ceb4f9892fa46bf738d08a0a4e25c885"
    end
  end

  def install
    bin.install "clwatch"
  end

  test do
    assert_match "clwatch #{version}", shell_output("#{bin}/clwatch version")
  end
end
