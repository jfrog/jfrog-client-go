const tasks = require("jfrog-pipelines-tasks");

async function downloadGo(version, targetFolder, cacheIntegration, cacheRepository) {
  const goUrl = computeDownloadUrl(version);
  tasks.info(`Go package url: ${goUrl}`);
  if (!cacheIntegration || !cacheRepository) {
    tasks.warning("Cache configuration not set. Caching will be skipped.");
  }
  const pathToFile = await tasks.downloadFile(goUrl, targetFolder, cacheRepository, cacheIntegration);
  await extractPackage(pathToFile, targetFolder);
}

function computeDownloadUrl(version) {
  const osFamily = getOsFamily();
  const architecture = getArchitecture();
  const packageType = getPackageType(osFamily);
  return `https://go.dev/dl/go${version}.${osFamily}-${architecture}.${packageType}`;
}

function getOsFamily() {
  return tasks.getOperatingSystemFamily().toLowerCase();
}

function getArchitecture() {
  const arch = tasks.getArchitecture();
  if (arch === "x86_64") return "amd64";
  if (arch === "ARM64") return "arm64";
  throw new tasks.PipelinesTaskError("Architecture not supported");
}

function getPackageType(osFamily) {
  if (osFamily === "windows") return "zip";
  return "tar.gz";
}

async function extractPackage(pathToPackage, targetFolder) {
  tasks.info("Extracting package content");
  const osFamily = getOsFamily();
  const packageType = getPackageType(osFamily);
  if (packageType === "zip")
    await tasks.unzip(pathToPackage, targetFolder);
  else
    await tasks.untar(pathToPackage, targetFolder);
}

module.exports = {
  downloadGo,
};
