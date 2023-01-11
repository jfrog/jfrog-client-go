const tasks = require("jfrog-pipelines-tasks");
const semver = require("semver");
const path = require("path");
const fs = require("fs");
const os = require("os");
const goDownloader = require("./goDownloader");

let inputVersion;
let inputCacheIntegration;
let inputCacheRepository;

function readAndValidateInput() {
  inputVersion = tasks.getInput("version");
  if (!inputVersion) throw "version input is required";

  inputVersion = semver.valid(inputVersion);
  if (!inputVersion) throw "version input must be semver compatible";

  inputCacheIntegration = tasks.getInput("cacheIntegration");
  inputCacheRepository = tasks.getInput("cacheRepository");
}

function findArtifactoryIntegration(inputCacheIntegration) {
  if (inputCacheIntegration) {
    const currentIntegration = tasks.getIntegration(inputCacheIntegration);
    if (currentIntegration.masterName.toLowerCase() === "artifactory") {
      return currentIntegration;
    } else {
      throw (
        "Input cacheIntegration is not an Artifactory Integration. Type: " +
        currentIntegration.masterName
      );
    }
  } else {
    tasks.info("Searching for Artifactory integration");
    try {
      const currentIntegration = tasks.findIntegrationByType("artifactory");
      tasks.info(`Artifactory integration ${currentIntegration.name} founded!`);
      return currentIntegration;
    } catch (err) {
      if (err instanceof tasks.IntegrationNotFound) {
        return undefined;
      } else {
        throw err;
      }
    }
  }
}

function createTargetFolder() {
  const goFolder = path.join(tasks.getStepWorkspaceDir(), "go");
  fs.mkdirSync(goFolder, { recursive: true });
  return goFolder;
}

async function setupEnvironment(targetFolder) {
  const goRoot = path.join(targetFolder, "go");
  const goBin = path.join(goRoot, "bin");

  tasks.info("Exporting GOROOT=" + goRoot);
  tasks.exportEnvironmentVariable("GOROOT", goRoot);

  tasks.info("Appending Go binaries location to PATH");
  tasks.appendToPath(goBin);
  const goPath = (await tasks.execute("go env GOPATH")).stdOut;
  if (goPath) {
    tasks.exportEnvironmentVariable("GOPATH", goPath);
    tasks.info("Appending GOPATH binaries location to PATH");
    const goPathBin = path.join(goPath, "bin");
    tasks.appendToPath(goPathBin);
  }
}

async function logGoEnvironment() {
  const goEnvOutput = (await tasks.execute("go env")).stdOut;
  tasks.info("Go env:" + os.EOL + goEnvOutput);
}

function logErrorAndExit(error) {
  tasks.error(error);
  tasks.debug(error.stack);
  process.exit(1);
}

async function run() {
  try {
    readAndValidateInput();
    const artifactoryIntegration = findArtifactoryIntegration(
      inputCacheIntegration
    );
    const targetFolder = createTargetFolder();
    await goDownloader.downloadGo(
      inputVersion,
      targetFolder,
      artifactoryIntegration,
      inputCacheRepository
    );
    await setupEnvironment(targetFolder);
    await logGoEnvironment();
  } catch (e) {
    logErrorAndExit(e);
  }
}

module.exports = {
  run,
  readAndValidateInput,
  findArtifactoryIntegration,
  createTargetFolder,
  setupEnvironment,
  logGoEnvironment,
  logErrorAndExit,
};
