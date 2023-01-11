const goDownloader = require("../goDownloader");
const tasks = require("jfrog-pipelines-tasks");
const utils = require("./__utils__/utils");
const path = require("path");
const fs = require("fs");

describe("go downloader tests", () => {
  afterEach(() => {
    jest.clearAllMocks();
    jest.resetAllMocks();
  });

  it("go downloader for x86_64 test", async () => {
    const tmpFolder = utils.createTempFolder("setup-go");
    fs.copyFileSync(
      path.join(__dirname, "__resources__", "go1.0.0.tgz"),
      path.join(tmpFolder, "go1.0.0.tgz")
    );
    const version = "1.0.0";
    jest.spyOn(tasks, "getOperatingSystemFamily").mockReturnValue("Linux");
    jest.spyOn(tasks, "getArchitecture").mockReturnValue("x86_64");
    const mockDownloadFile = jest
      .spyOn(tasks, "downloadFile")
      .mockResolvedValue(path.join(tmpFolder, "go1.0.0.tgz"));

    await goDownloader.downloadGo(version, tmpFolder);

    expect(mockDownloadFile).toHaveBeenCalledWith(
      "https://go.dev/dl/go1.0.0.linux-amd64.tar.gz",
      tmpFolder,
      undefined,
      undefined
    );
    expect(
      fs.existsSync(path.join(tmpFolder, "go", "go.version"))
    ).toBeTruthy();
  });

  it("go downloader for arm test", async () => {
    const tmpFolder = utils.createTempFolder("setup-go");
    fs.copyFileSync(
      path.join(__dirname, "__resources__", "go1.0.0.tgz"),
      path.join(tmpFolder, "go1.0.0.tgz")
    );
    const version = "1.0.0";
    jest.spyOn(tasks, "getOperatingSystemFamily").mockReturnValue("Linux");
    jest.spyOn(tasks, "getArchitecture").mockReturnValue("ARM64");
    const mockDownloadFile = jest
      .spyOn(tasks, "downloadFile")
      .mockResolvedValue(path.join(tmpFolder, "go1.0.0.tgz"));

    await goDownloader.downloadGo(version, tmpFolder);

    expect(mockDownloadFile).toHaveBeenCalledWith(
      "https://go.dev/dl/go1.0.0.linux-arm64.tar.gz",
      tmpFolder,
      undefined,
      undefined
    );

    expect(
      fs.existsSync(path.join(tmpFolder, "go", "go.version"))
    ).toBeTruthy();
  });

  it("arch not supported error", async () => {
    const tmpFolder = utils.createTempFolder("setup-go");
    fs.copyFileSync(
      path.join(__dirname, "__resources__", "go1.0.0.tgz"),
      path.join(tmpFolder, "go1.0.0.tgz")
    );
    const version = "1.0.0";
    jest.spyOn(tasks, "getOperatingSystemFamily").mockReturnValue("Linux");
    jest.spyOn(tasks, "getArchitecture").mockReturnValue("NOT_SUPPORTED");
    const mockDownloadFile = jest
      .spyOn(tasks, "downloadFile")
      .mockResolvedValue(path.join(tmpFolder, "go1.0.0.tgz"));

    expect.assertions(1);
    try {
      await goDownloader.downloadGo(version, tmpFolder);
    } catch (err) {
      expect(err).toBeInstanceOf(tasks.PipelinesTaskError);
    }
  });

  it("go downloader for windows test", async () => {
    const tmpFolder = utils.createTempFolder("setup-go");
    fs.copyFileSync(
      path.join(__dirname, "__resources__", "go1.0.0.zip"),
      path.join(tmpFolder, "go1.0.0.zip")
    );
    const version = "1.0.0";
    jest.spyOn(tasks, "getOperatingSystemFamily").mockReturnValue("Windows");
    jest.spyOn(tasks, "getArchitecture").mockReturnValue("x86_64");
    const mockDownloadFile = jest
      .spyOn(tasks, "downloadFile")
      .mockResolvedValue(path.join(tmpFolder, "go1.0.0.zip"));

    await goDownloader.downloadGo(version, tmpFolder);

    expect(mockDownloadFile).toHaveBeenCalledWith(
      "https://go.dev/dl/go1.0.0.windows-amd64.zip",
      tmpFolder,
      undefined,
      undefined
    );

    expect(
      fs.existsSync(path.join(tmpFolder, "go", "go.version"))
    ).toBeTruthy();
  });
});
