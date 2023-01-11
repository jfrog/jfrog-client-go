const setupGo = require("../setupGo");
const tasks = require("jfrog-pipelines-tasks");

describe("input validation tests", () => {
  afterEach(() => {
    jest.clearAllMocks();
    jest.resetAllMocks();
  });

  it("no version provided should throw error", () => {
    expect(setupGo.readAndValidateInput).toThrow("version input is required");
  });

  it("invalid version should throw error", () => {
    jest.spyOn(tasks, "getInput").mockImplementation((name) => {
      if (name === "version") return "a.b.c";
    });
    expect(setupGo.readAndValidateInput).toThrow(
      "version input must be semver compatible"
    );
  });

  it("valid version should succeed", () => {
    jest.spyOn(tasks, "getInput").mockImplementation((name) => {
      if (name === "version") return "0.0.1";
    });
    setupGo.readAndValidateInput();
  });
});
