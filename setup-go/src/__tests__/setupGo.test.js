const setupGo = require("../setupGo");
const utils = require("./__utils__/utils");
const tasks = require("jfrog-pipelines-tasks");
const goDownloader = require("../goDownloader");
const fs = require("fs");
const path = require("path");

describe("log error and exit tests", () => {
  afterEach(() => {
    jest.clearAllMocks();
    jest.resetAllMocks();
  });

  it("Create target folder", () => {
    const tmpFolder = utils.createTempFolder("setup-go");
    jest.spyOn(tasks, "getStepWorkspaceDir").mockReturnValue(tmpFolder);
    setupGo.createTargetFolder();
    const expectedFolder = path.join(tmpFolder, "go");
    expect(fs.existsSync(expectedFolder)).toBeTruthy();
  });

  it("log error and exit", () => {
    const mockExit = jest.spyOn(process, "exit").mockImplementation(() => {});
    setupGo.logErrorAndExit("error");
    expect(mockExit).toHaveBeenCalledWith(1);
  });

  it("run task", async () => {
    const tmpFolder = utils.createTempFolder("setup-go");
    jest.spyOn(tasks, "getStepWorkspaceDir").mockReturnValue(tmpFolder);
    jest.spyOn(tasks, "getInput").mockImplementation((name) => {
      if (name === "version") return "0.0.1";
    });
    jest.spyOn(tasks, "findIntegrationByType").mockReturnValue({
      id: 1,
      name: "artifactory",
      masterName: "Artifactory",
      displayName: "Artifactory",
    });
    jest.spyOn(goDownloader, "downloadGo").mockResolvedValue(null);
    jest.spyOn(tasks, "exportEnvironmentVariable").mockReturnThis();
    jest.spyOn(tasks, "appendToPath").mockReturnThis();
    jest.spyOn(tasks, "execute").mockResolvedValue({
      stdOut: "",
    });
    await setupGo.run();
  });

  it("run task with error", async () => {
    const mockExit = jest.spyOn(process, "exit").mockImplementation(() => {});
    await setupGo.run();
    expect(mockExit).toHaveBeenCalledWith(1);
  });
});
