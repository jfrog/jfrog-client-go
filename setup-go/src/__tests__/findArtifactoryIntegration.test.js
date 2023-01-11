const setupGo = require("../setupGo");
const tasks = require("jfrog-pipelines-tasks");

describe("find Artifactory Integration tests", () => {
  afterEach(() => {
    jest.clearAllMocks();
    jest.resetAllMocks();
  });

  it("with valid integration input", () => {
    const getIntegrationMock = jest
      .spyOn(tasks, "getIntegration")
      .mockImplementation((name) => {
        return {
          id: 1,
          name: name,
          masterName: "Artifactory",
          displayName: "Artifactory",
        };
      });
    setupGo.findArtifactoryIntegration("integration");
    expect(getIntegrationMock).toHaveBeenCalledWith("integration");
  });

  it("with invalid integration input", () => {
    const getIntegrationMock = jest
      .spyOn(tasks, "getIntegration")
      .mockImplementation((name) => {
        return {
          id: 1,
          name: name,
          masterName: "GitHub",
          displayName: "GitHub",
        };
      });
    expect(
      setupGo.findArtifactoryIntegration.bind(null, "integration")
    ).toThrow(
      "Input cacheIntegration is not an Artifactory Integration. Type: GitHub"
    );
    expect(getIntegrationMock).toHaveBeenCalledWith("integration");
  });

  it("with artifactory integration available", () => {
    const findIntegrationMock = jest
      .spyOn(tasks, "findIntegrationByType")
      .mockReturnValue({
        id: 1,
        name: "artifactory",
        masterName: "Artifactory",
        displayName: "Artifactory",
      });
    setupGo.findArtifactoryIntegration("");
    expect(findIntegrationMock).toHaveBeenCalledWith("artifactory");
  });

  it("artifactory integration not available", () => {
    const findIntegrationMock = jest
      .spyOn(tasks, "findIntegrationByType")
      .mockImplementation((type) => {
        throw new tasks.IntegrationNotFound(type);
      });

    const artIntegration = setupGo.findArtifactoryIntegration("");
    expect(artIntegration).toBeUndefined();
    expect(findIntegrationMock).toHaveBeenCalledWith("artifactory");
  });

  it("unexpected error should throw exception", () => {
    const findIntegrationMock = jest
      .spyOn(tasks, "findIntegrationByType")
      .mockImplementation((type) => {
        throw "unexpected error";
      });

    expect(setupGo.findArtifactoryIntegration.bind(null, "")).toThrow(
      "unexpected error"
    );
    expect(findIntegrationMock).toHaveBeenCalledWith("artifactory");
  });
});
