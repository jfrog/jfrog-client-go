const path = require("path");
const os = require("os");
const {randomUUID} = require("crypto");
const fs = require("fs");

function createTempFolder(prefix) {
  const tmpFolder = path.join(os.tmpdir(), prefix + "-" + randomUUID());
  fs.mkdirSync(tmpFolder);
  return tmpFolder;
}

module.exports = {
  createTempFolder
}