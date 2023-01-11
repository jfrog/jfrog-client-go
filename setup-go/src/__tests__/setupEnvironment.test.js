const setupGo = require("../setupGo");
const utils = require("./__utils__/utils");
const tasks = require("jfrog-pipelines-tasks");
const path = require("path");
const os = require("os");

describe("setup environment tests", () => {
  afterEach(() => {
    jest.clearAllMocks();
    jest.resetAllMocks();
  });

  it("setup environment test", async () => {
    const mockGoPath = "/home/user/gopath";
    const exportedEnvVars = {};
    jest
      .spyOn(tasks, "exportEnvironmentVariable")
      .mockImplementation((key, value) => {
        exportedEnvVars[key] = value;
      });

    const appendedToPath = [];
    jest
      .spyOn(tasks, "appendToPath")
      .mockImplementation((value) => {
        appendedToPath.push(value);
      });
    jest.spyOn(tasks, "execute").mockResolvedValue({
      stdOut: mockGoPath
    })

    const goFolder = utils.createTempFolder("setup-go");

    const expectedExportedEnvVars = {
      GOROOT: path.join(goFolder, "go"),
      GOPATH: mockGoPath
    }
    const expectedAppendToPath = [
      path.join(goFolder, "go", "bin"),
      path.join(mockGoPath, "bin")
    ]

    await setupGo.setupEnvironment(goFolder);

    expect(exportedEnvVars).toEqual(expectedExportedEnvVars);
    expect(appendedToPath).toEqual(expectedAppendToPath);
  });

  it("log go environment", async () => {
    const execMock = jest.spyOn(tasks, "execute").mockResolvedValue({
      stdOut: "output"
    })
    const logMock = jest.spyOn(tasks, "info").mockReturnThis();
    await setupGo.logGoEnvironment();
    expect(execMock).toHaveBeenCalledWith("go env");
    expect(logMock).toHaveBeenCalledWith("Go env:" + os.EOL + "output");
  });
});
