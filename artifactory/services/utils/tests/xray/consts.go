package xray

import (
	xscServices "github.com/jfrog/jfrog-client-go/xsc/services"
)

const ScanResponse = `
{
"scan_id": "3472b4e2-bddc-11ee-a9c9-acde48001122",
	"vulnerabilities": [{
		"summary": "test",
		"severity": "high"
	}]
}
`
const FatalErrorXrayScanResponse = `
{
 "errors": [{"status":-1}, {"status":500}]
}
`

const CleanXrayScanResponse = `
{
  "summary" : {
    "message" : "Build pipeline-xray number 307 was scanned by Xray and passed with no Alerts",
    "total_alerts" : 0,
    "fail_build" : false,
    "more_details_url" : ""
  },
  "alerts" : [ ],
  "licenses" : [ {
    "name" : "Unknown",
    "components" : [ "generic://sha256:8739c76e681f900923b900c9df0ef75cf421d39cabb54650c4b9ad19b6a76d85/ArtifactoryPipeline.zip", "gav://com.opensymphony:xwork:2.0.6", "build://pipeline-xray:307" ],
    "full_name" : "Unknown license",
    "more_info_url" : null
  } ]
}`

const VulnerableXrayScanResponse = `{
  "summary": {
    "message": "Build pipeline-xray number 294 was scanned by Xray and 2 Alerts were generated",
    "total_alerts": 2,
    "fail_build": true,
    "more_details_url": "http://10.1.20.29:8000//web/#/alerts/all?filter=58760987e138231efd444492,58760987e138231efd444493"
  },
  "alerts": [
    {
      "created": "2017-01-11T12:31:35.79130964+02:00",
      "issues": [
        {
          "severity": "Major",
          "type": "security",
          "provider": "JFrog",
          "created": "2011-05-13T16:05:45+03:00",
          "summary": "CWE-200 Information Exposure",
          "description": "XWork 2.2.1 in Apache Struts 2.2.1, and OpenSymphony XWork in OpenSymphony WebWork, allows remote attackers to obtain potentially sensitive information about internal Java class paths via vectors involving an s:submit element and a nonexistent method, a different vulnerability than CVE-2011-1772.3.",
          "cve": "CVE-2011-2088",
          "component_ids": [
            "gav://opensymphony:xwork:1.0.3.6",
            "gav://opensymphony:xwork:1.0.3.5",
            "gav://opensymphony:xwork:1.0.3.4",
            "gav://opensymphony:xwork:1.0.3.3",
            "gav://opensymphony:xwork:1.1.1",
            "gav://opensymphony:xwork:1.0.3.2",
            "gav://opensymphony:xwork:1.0.3.1",
            "gav://opensymphony:xwork:1.1.3",
            "gav://opensymphony:webwork:2.1.5-atlassian-2",
            "gav://opensymphony:webwork:2.1.5-atlassian-3",
            "gav://opensymphony:webwork:2.2.1",
            "gav://opensymphony:webwork:2.2.2",
            "gav://opensymphony:webwork:2.2.3",
            "gav://opensymphony:xwork:2.0-beta-3-j4",
            "gav://opensymphony:xwork:2.0.1-j4",
            "gav://com.opensymphony:webwork:2.2.6-atlassian-3",
            "gav://opensymphony:webwork:2.2.4",
            "gav://opensymphony:webwork:2.2.5",
            "gav://com.opensymphony:webwork:2.2.6-atlassian-2",
            "gav://opensymphony:webwork:2.2.7",
            "gav://opensymphony:xwork:2.0.3",
            "gav://opensymphony:xwork:2.0.2",
            "gav://com.opensymphony:xwork:2.0.7",
            "gav://com.opensymphony:xwork:2.0.5",
            "gav://com.opensymphony:xwork:2.0.6",
            "gav://opensymphony:webwork:2.2.7-atlassian-25",
            "gav://opensymphony:webwork:2.2.7-atlassian-27",
            "gav://opensymphony:xwork:2.0.0-j4",
            "gav://opensymphony:xwork:1.2.2",
            "gav://opensymphony:xwork:1.2.1",
            "gav://opensymphony:webwork:2.2.7-atlassian-23",
            "gav://opensymphony:xwork:1.2.3",
            "gav://opensymphony:webwork:2.2.7-atlassian-24",
            "gav://opensymphony:xwork:2.0-beta-2",
            "gav://opensymphony:xwork:2.0-RC1",
            "gav://opensymphony:xwork:2.0-beta-3",
            "gav://opensymphony:xwork:2.0-RC1-j4",
            "gav://opensymphony:xwork:2.0-beta-1",
            "gav://opensymphony:xwork:2.0.1",
            "gav://opensymphony:xwork:2.0.0",
            "gav://opensymphony:xwork:1.1",
            "gav://opensymphony:xwork:1.2",
            "gav://opensymphony:xwork:1.0",
            "gav://opensymphony:webwork:2.2.7-atlassian-29",
            "gav://opensymphony:webwork:1.4-atlassian-26-testmvn2-1",
            "gav://opensymphony:webwork:2.2.7-atlassian-1",
            "gav://opensymphony:webwork:2.2.7-atlassian-2",
            "gav://opensymphony:webwork:2.2.7-atlassian-3",
            "gav://opensymphony:webwork:2.2.7-atlassian-5",
            "gav://opensymphony:webwork:2.2.7-atlassian-6",
            "gav://com.opensymphony:webwork:2.2.6",
            "gav://opensymphony:webwork:1.4-atlassian-9",
            "gav://opensymphony:webwork:2.1",
            "gav://opensymphony:webwork:1.4-atlassian-4",
            "gav://opensymphony:webwork:2.0",
            "gav://opensymphony:webwork:1.4-atlassian-2",
            "gav://opensymphony:webwork:2.2",
            "gav://opensymphony:webwork:1.4-atlassian-1",
            "gav://com.opensymphony:xwork:2.0.4",
            "gav://opensymphony:xwork:1.2.5-rc1",
            "gav://opensymphony:xwork:1.0.1",
            "gav://opensymphony:webwork:1.4-atlassian-19",
            "gav://opensymphony:xwork:1.2.5-atlassian-1",
            "gav://opensymphony:xwork:1.2.5-atlassian-2",
            "gav://opensymphony:xwork:1.2.5-atlassian-4",
            "gav://opensymphony:xwork:1.0.4",
            "gav://opensymphony:xwork:1.0.3",
            "gav://opensymphony:xwork:1.2.5-atlassian-5",
            "gav://opensymphony:xwork:1.2.5-atlassian-6",
            "gav://opensymphony:xwork:1.2.5-atlassian-7",
            "gav://opensymphony:xwork:1.0.5",
            "gav://opensymphony:webwork:1.4-atlassian-10",
            "gav://opensymphony:webwork:1.4-atlassian-11",
            "gav://opensymphony:webwork:1.4-atlassian-12",
            "gav://opensymphony:webwork:1.4-atlassian-13",
            "gav://opensymphony:webwork:1.4-atlassian-15",
            "gav://opensymphony:webwork:1.4-atlassian-16",
            "gav://opensymphony:webwork:1.4-atlassian-17",
            "gav://opensymphony:xwork:1.2.3-20070717",
            "gav://opensymphony:webwork:1.4-atlassian-30",
            "gav://opensymphony:webwork:12Dec05-jiratld",
            "gav://opensymphony:xwork:1.2.5-atlassian-8",
            "gav://opensymphony:webwork:1.4",
            "gav://opensymphony:webwork:2.1.3",
            "gav://opensymphony:webwork:2.1.4",
            "gav://com.opensymphony:xwork:2.1.2",
            "gav://opensymphony:webwork:2.1.5",
            "gav://com.opensymphony:xwork:2.1.3",
            "gav://opensymphony:webwork:2.1.6",
            "gav://com.opensymphony:xwork:2.1.0",
            "gav://opensymphony:webwork:2.1.7",
            "gav://com.opensymphony:xwork:2.1.1",
            "gav://opensymphony:webwork:1.4-atlassian-22",
            "gav://opensymphony:webwork:1.4-atlassian-23",
            "gav://opensymphony:webwork:1.4-atlassian-24",
            "gav://opensymphony:webwork:1.4-atlassian-25",
            "gav://opensymphony:webwork:1.4-atlassian-27",
            "gav://opensymphony:xwork-src:1.1.3",
            "gav://opensymphony:xwork-tiger-src:1.2.2",
            "gav://opensymphony:xwork-tiger-src:1.2.3",
            "gav://opensymphony:xwork-tiger-src:1.2.1",
            "gav://org.apache.struts:struts2-embeddedjsp-plugin:2.2.1",
            "gav://opensymphony:xwork-src:1.1.1",
            "gav://org.apache.struts:struts2-archetype-blank:2.2.1",
            "gav://org.apache.struts:struts2-portlet:2.2.1",
            "gav://org.apache.struts:struts2-archetype-starter:2.2.1",
            "gav://org.apache.struts:struts2-spring-plugin:2.2.1",
            "gav://org.apache.struts:struts2-jasperreports-plugin:2.2.1",
            "gav://com.opensymphony:xwork-parent:2.1.6",
            "gav://com.opensymphony:xwork-parent:2.1.5",
            "gav://org.apache.struts:struts2-archetypes:2.2.1",
            "gav://org.apache.struts:struts2-assembly:2.2.1",
            "gav://com.opensymphony:xwork-parent:2.1.4",
            "gav://org.apache.struts:struts2-jboss-blank:2.2.1",
            "gav://org.apache.struts:struts2-sitegraph-plugin:2.2.1",
            "gav://org.apache.struts:struts2-showcase:2.2.1",
            "gav://com.opensymphony:xwork-core:2.1.5",
            "gav://org.apache.struts:struts2-core:2.2.1",
            "gav://com.opensymphony:xwork-core:2.1.4",
            "gav://org.apache.struts:struts2-portlet-plugin:2.2.1",
            "gav://com.opensymphony:xwork-core:2.1.6",
            "gav://org.apache.struts:struts2-osgi-demo-bundle:2.2.1",
            "gav://com.opensymphony.webwork:com.springsource.com.opensymphony.webwork:2.2.6",
            "gav://com.opensymphony.xwork:com.springsource.com.opensymphony.xwork:1.2.3",
            "gav://opensymphony:xwork-tiger-src:1.1.3",
            "gav://opensymphony:xwork-src:1.2.2",
            "gav://opensymphony:xwork-src:1.2.3",
            "gav://opensymphony:xwork-tiger-src:1.1.1",
            "gav://org.apache.struts:struts2-tiles-plugin:2.2.1",
            "gav://org.apache.struts:struts2-junit-plugin:2.2.1",
            "gav://opensymphony:webwork-src:2.2",
            "gav://opensymphony:xwork-src:1.2.1",
            "gav://org.apache.struts:struts2-dojo-plugin:2.2.1",
            "gav://org.apache.struts:struts2-json-plugin:2.2.1",
            "gav://org.apache.struts:struts2-osgi-bundles:2.2.1",
            "gav://org.apache.struts:struts2-osgi-admin-bundle:2.2.1",
            "gav://opensymphony:xwork-src:2.0-beta-1",
            "gav://opensymphony:xwork-src:2.0-beta-2",
            "gav://opensymphony:webwork-src:2.2.3",
            "gav://opensymphony:webwork-src:2.2.4",
            "gav://opensymphony:webwork-src:2.2.1",
            "gav://org.apache.struts:struts2-dwr-plugin:2.2.1",
            "gav://opensymphony:webwork-src:2.2.2",
            "gav://opensymphony:xwork-tiger:1.2",
            "gav://org.apache.struts:struts2-oval-plugin:2.2.1",
            "gav://opensymphony:xwork-tiger:1.1",
            "gav://org.apache.struts:struts2-codebehind-plugin:2.2.1",
            "gav://opensymphony:webwork-src:2.2.5",
            "gav://opensymphony:xwork-tiger:1.2.3",
            "gav://opensymphony:xwork-tiger:1.2.1",
            "gav://com.opensymphony.xwork:com.springsource.com.opensymphony.xwork:1.1.1",
            "gav://opensymphony:xwork-tiger:1.2.2",
            "gav://org.apache.struts:struts2-archetype-dbportlet:2.2.1",
            "gav://opensymphony:xwork-tiger-src:1.2",
            "gav://org.apache.struts:struts2-plexus-plugin:2.2.1",
            "gav://opensymphony:xwork-tiger-src:1.1",
            "gav://org.apache.struts:struts2-jfreechart-plugin:2.2.1",
            "gav://org.apache.struts:struts2-jsf-plugin:2.2.1",
            "gav://opensymphony:xwork-src:1.1",
            "gav://org.apache.struts:struts2-pell-multipart-plugin:2.2.1",
            "gav://opensymphony:xwork-src:1.2",
            "gav://com.opensymphony:xwork-assembly:2.1.4",
            "gav://org.apache.struts:struts2-rest-showcase:2.2.1",
            "gav://com.opensymphony:xwork-assembly:2.1.6",
            "gav://org.apache.struts:struts2-struts1-plugin:2.2.1",
            "gav://org.apache.struts:struts2-javatemplates-plugin:2.2.1",
            "gav://com.opensymphony:xwork-assembly:2.1.5",
            "gav://org.apache.struts:struts2-osgi-plugin:2.2.1",
            "gav://org.apache.struts:struts2-archetype-convention:2.2.1",
            "gav://opensymphony:xwork-tiger:2.0-beta-1",
            "gav://opensymphony:xwork-tiger-src:2.0-beta-1",
            "gav://org.apache.struts:struts2-archetype-plugin:2.2.1",
            "gav://opensymphony:xwork-tiger:1.1.3",
            "gav://opensymphony:xwork-tiger:1.1.1",
            "gav://org.apache.struts:struts2-blank:2.2.1",
            "gav://org.apache.struts:struts2-testng-plugin:2.2.1",
            "gav://org.apache.struts:struts2-parent:2.2.1",
            "gav://org.apache.struts:struts2-sitemesh-plugin:2.2.1",
            "gav://org.apache.struts:struts2-config-browser-plugin:2.2.1",
            "gav://org.apache.struts:struts2-convention-plugin:2.2.1",
            "gav://org.apache.struts:struts2-plugins:2.2.1",
            "gav://org.apache.struts:struts2-rest-plugin:2.2.1",
            "gav://org.apache.struts:struts2-apps:2.2.1",
            "gav://org.apache.struts.xwork:xwork-core:2.2.1",
            "gav://org.apache.struts:struts2-mailreader:2.2.1",
            "gav://com.opensymphony:xwork-plugins:2.1.4",
            "gav://com.opensymphony:xwork-plugins:2.1.5",
            "gav://com.opensymphony:xwork-plugins:2.1.6",
            "gav://org.apache.struts:struts2-archetype-portlet:2.2.1"
          ],
          "impacted_artifacts": [
            {
              "name": "pipeline-xray",
              "path": "LOCAL/builds/",
              "sha256": "f8fa924e200c8412075de89d867f1b8dd5aa343938754fb652f84ff472a75b97",
              "sha1": "",
              "depth": 0,
              "display_name": "pipeline-xray:294",
              "pkg_type": "Build",
              "parent_sha": "f8fa924e200c8412075de89d867f1b8dd5aa343938754fb652f84ff472a75b97",
              "impact_path": "",
              "infected_file": {
                "name": "xwork-2.0.6.jar",
                "path": "LOCAL/libs-release-local/com/opensymphony/xwork/2.0.6/",
                "sha256": "4baa7d0a7203ceeefa063a3782de733d417c4eae5ee5c0ab06e9e9e4607c004a",
                "sha1": "",
                "depth": 0,
                "parent_sha": "4baa7d0a7203ceeefa063a3782de733d417c4eae5ee5c0ab06e9e9e4607c004a",
                "display_name": "com.opensymphony:xwork:2.0.6",
                "pkg_type": "Maven"
              }
            }
          ]
        },
        {
          "severity": "Minor",
          "type": "security",
          "provider": "JFrog",
          "created": "2011-05-13T16:05:44+03:00",
          "summary": "CWE-79 Improper Neutralization of Input During Web Page Generation ('Cross-site Scripting')",
          "description": "Multiple cross-site scripting (XSS) vulnerabilities in XWork in Apache Struts 2.x before 2.2.3, and OpenSymphony XWork in OpenSymphony WebWork, allow remote attackers to inject arbitrary web script or HTML via vectors involving (1) an action name, (2) the action attribute of an s:submit element, or (3) the method attribute of an s:submit element.",
          "cve": "CVE-2011-1772",
          "component_ids": [
            "gav://opensymphony:xwork:1.0.3.6",
            "gav://opensymphony:xwork:1.0.3.5",
            "gav://opensymphony:xwork:1.0.3.4",
            "gav://opensymphony:xwork:1.0.3.3",
            "gav://opensymphony:xwork:1.1.1",
            "gav://opensymphony:xwork:1.0.3.2",
            "gav://opensymphony:xwork:1.0.3.1",
            "gav://opensymphony:xwork:1.1.3",
            "gav://opensymphony:webwork:2.1.5-atlassian-2",
            "gav://opensymphony:webwork:2.1.5-atlassian-3",
            "gav://opensymphony:webwork:2.2.1",
            "gav://opensymphony:webwork:2.2.2",
            "gav://opensymphony:webwork:2.2.3",
            "gav://opensymphony:xwork:2.0-beta-3-j4",
            "gav://opensymphony:xwork:2.0.1-j4",
            "gav://com.opensymphony:webwork:2.2.6-atlassian-3",
            "gav://opensymphony:webwork:2.2.4",
            "gav://opensymphony:webwork:2.2.5",
            "gav://com.opensymphony:webwork:2.2.6-atlassian-2",
            "gav://opensymphony:webwork:2.2.7",
            "gav://opensymphony:xwork:2.0.3",
            "gav://opensymphony:xwork:2.0.2",
            "gav://com.opensymphony:xwork:2.0.7",
            "gav://com.opensymphony:xwork:2.0.5",
            "gav://com.opensymphony:xwork:2.0.6",
            "gav://opensymphony:webwork:2.2.7-atlassian-25",
            "gav://opensymphony:webwork:2.2.7-atlassian-27",
            "gav://opensymphony:xwork:2.0.0-j4",
            "gav://opensymphony:xwork:1.2.2",
            "gav://opensymphony:xwork:1.2.1",
            "gav://opensymphony:webwork:2.2.7-atlassian-23",
            "gav://opensymphony:xwork:1.2.3",
            "gav://opensymphony:webwork:2.2.7-atlassian-24",
            "gav://opensymphony:xwork:2.0-beta-2",
            "gav://opensymphony:xwork:2.0-RC1",
            "gav://opensymphony:xwork:2.0-beta-3",
            "gav://opensymphony:xwork:2.0-RC1-j4",
            "gav://opensymphony:xwork:2.0-beta-1",
            "gav://opensymphony:xwork:2.0.1",
            "gav://opensymphony:xwork:2.0.0",
            "gav://opensymphony:xwork:1.1",
            "gav://opensymphony:xwork:1.2",
            "gav://opensymphony:xwork:1.0",
            "gav://opensymphony:webwork:2.2.7-atlassian-29",
            "gav://opensymphony:webwork:1.4-atlassian-26-testmvn2-1",
            "gav://opensymphony:webwork:2.2.7-atlassian-1",
            "gav://opensymphony:webwork:2.2.7-atlassian-2",
            "gav://opensymphony:webwork:2.2.7-atlassian-3",
            "gav://opensymphony:webwork:2.2.7-atlassian-5",
            "gav://opensymphony:webwork:2.2.7-atlassian-6",
            "gav://com.opensymphony:webwork:2.2.6",
            "gav://opensymphony:webwork:1.4-atlassian-9",
            "gav://opensymphony:webwork:2.1",
            "gav://opensymphony:webwork:1.4-atlassian-4",
            "gav://opensymphony:webwork:2.0",
            "gav://opensymphony:webwork:1.4-atlassian-2",
            "gav://opensymphony:webwork:2.2",
            "gav://opensymphony:webwork:1.4-atlassian-1",
            "gav://com.opensymphony:xwork:2.0.4",
            "gav://opensymphony:xwork:1.2.5-rc1",
            "gav://opensymphony:xwork:1.0.1",
            "gav://opensymphony:webwork:1.4-atlassian-19",
            "gav://opensymphony:xwork:1.2.5-atlassian-1",
            "gav://opensymphony:xwork:1.2.5-atlassian-2",
            "gav://opensymphony:xwork:1.2.5-atlassian-4",
            "gav://opensymphony:xwork:1.0.4",
            "gav://opensymphony:xwork:1.0.3",
            "gav://opensymphony:xwork:1.2.5-atlassian-5",
            "gav://opensymphony:xwork:1.2.5-atlassian-6",
            "gav://opensymphony:xwork:1.2.5-atlassian-7",
            "gav://opensymphony:xwork:1.0.5",
            "gav://opensymphony:webwork:1.4-atlassian-10",
            "gav://opensymphony:webwork:1.4-atlassian-11",
            "gav://opensymphony:webwork:1.4-atlassian-12",
            "gav://opensymphony:webwork:1.4-atlassian-13",
            "gav://opensymphony:webwork:1.4-atlassian-15",
            "gav://opensymphony:webwork:1.4-atlassian-16",
            "gav://opensymphony:webwork:1.4-atlassian-17",
            "gav://opensymphony:xwork:1.2.3-20070717",
            "gav://opensymphony:webwork:1.4-atlassian-30",
            "gav://opensymphony:webwork:12Dec05-jiratld",
            "gav://opensymphony:xwork:1.2.5-atlassian-8",
            "gav://opensymphony:webwork:1.4",
            "gav://opensymphony:webwork:2.1.3",
            "gav://opensymphony:webwork:2.1.4",
            "gav://com.opensymphony:xwork:2.1.2",
            "gav://opensymphony:webwork:2.1.5",
            "gav://com.opensymphony:xwork:2.1.3",
            "gav://opensymphony:webwork:2.1.6",
            "gav://com.opensymphony:xwork:2.1.0",
            "gav://opensymphony:webwork:2.1.7",
            "gav://com.opensymphony:xwork:2.1.1",
            "gav://opensymphony:webwork:1.4-atlassian-22",
            "gav://opensymphony:webwork:1.4-atlassian-23",
            "gav://opensymphony:webwork:1.4-atlassian-24",
            "gav://opensymphony:webwork:1.4-atlassian-25",
            "gav://opensymphony:webwork:1.4-atlassian-27",
            "gav://org.apache.struts:struts2-testng-plugin:2.1.6",
            "gav://org.apache.struts:struts2-sitegraph-plugin:2.0.12",
            "gav://org.apache.struts:struts2-sitegraph-plugin:2.0.11",
            "gav://org.apache.struts:struts2-portlet:2.2.1.1",
            "gav://org.apache.struts:struts2-testng-plugin:2.1.8",
            "gav://org.apache.struts:struts2-jfreechart-plugin:2.0.8",
            "gav://org.apache.struts:struts2-testng-plugin:2.1.2",
            "gav://org.apache.struts:struts2-jfreechart-plugin:2.0.5",
            "gav://org.apache.struts:struts2-jfreechart-plugin:2.0.6",
            "gav://org.apache.struts:struts2-portlet:2.2.1",
            "gav://org.apache.struts:struts2-archetype-starter:2.2.1.1",
            "gav://org.apache.struts:struts2-spring-plugin:2.2.1",
            "gav://org.apache.struts:struts2-config-browser-plugin:2.1.6",
            "gav://org.apache.struts:struts2-config-browser-plugin:2.1.8",
            "gav://org.apache.struts:struts2-config-browser-plugin:2.1.2",
            "gav://org.apache.struts:struts2-tiles-plugin:2.1.8.1",
            "gav://org.apache.struts:struts2-jfreechart-plugin:2.1.8",
            "gav://org.apache.struts:struts2-jsf-plugin:2.1.8",
            "gav://org.apache.struts:struts2-jsf-plugin:2.1.6",
            "gav://org.apache.struts:struts2-dojo-plugin:2.1.8.1",
            "gav://org.apache.struts:struts2-archetype-convention:2.2.1.1",
            "gav://org.apache.struts:struts2-spring-plugin:2.2.1.1",
            "gav://org.apache.struts:struts2-jsf-plugin:2.0.11.1",
            "gav://org.apache.struts:struts2-jfreechart-plugin:2.1.2",
            "gav://org.apache.struts:struts2-jsf-plugin:2.0.11.2",
            "gav://org.apache.struts:struts2-codebehind-plugin:2.0.11.1",
            "gav://org.apache.struts:struts2-jfreechart-plugin:2.1.6",
            "gav://org.apache.struts:struts2-codebehind-plugin:2.0.11.2",
            "gav://org.apache.struts:struts2-config-browser-plugin:2.0.6",
            "gav://org.apache.struts:struts2-config-browser-plugin:2.0.9",
            "gav://org.apache.struts:struts2-config-browser-plugin:2.0.8",
            "gav://org.apache.struts:struts2-jboss-blank:2.2.1.1",
            "gav://org.apache.struts:struts2-config-browser-plugin:2.0.5",
            "gav://org.apache.struts:struts2-dwr-plugin:2.2.1",
            "gav://org.apache.struts:struts-parent:2",
            "gav://org.apache.struts:struts2-oval-plugin:2.2.1",
            "gav://org.apache.struts:struts2-osgi-demo-bundle:2.2.1.1",
            "gav://org.apache.struts:struts2-jsf-plugin:2.0.8",
            "gav://org.apache.struts:struts2-jasperreports-plugin:2.2.1.1",
            "gav://org.apache.struts:struts2-jsf-plugin:2.0.9",
            "gav://org.apache.struts:struts2-jsf-plugin:2.0.6",
            "gav://org.apache.struts:struts2-testng-plugin:2.2.1.1",
            "gav://org.apache.struts:struts2-jsf-plugin:2.0.5",
            "gav://org.apache.struts:struts2-jsf-plugin:2.2.1.1",
            "gav://org.apache.struts:struts2-rest-plugin:2.1.8.1",
            "gav://org.apache.struts:struts2-plexus-plugin:2.1.8.1",
            "gav://org.apache.struts:struts2-jfreechart-plugin:2.2.1",
            "gav://org.apache.struts:struts2-config-browser-plugin:2.2.1.1",
            "gav://org.apache.struts:struts2-archetype-blank:2.2.1.1",
            "gav://org.apache.struts:struts2-archetype-convention:2.2.1",
            "gav://org.apache.struts:struts2-oval-plugin:2.1.8",
            "gav://opensymphony:xwork-tiger:2.0-beta-1",
            "gav://org.apache.struts:struts2-osgi-admin-bundle:2.1.8.1",
            "gav://org.apache.struts:struts2-rest-showcase:2.2.1.1",
            "gav://org.apache.struts:struts2-config-browser-plugin:2.0.14",
            "gav://org.apache.empire-db:empire-db-struts2:2.0.6-incubating",
            "gav://org.apache.struts:struts2-parent:2.2.1",
            "gav://org.apache.struts:struts2-config-browser-plugin:2.0.11",
            "gav://org.apache.struts:struts2-config-browser-plugin:2.0.12",
            "gav://org.apache.struts:struts2-pell-multipart-plugin:2.0.12",
            "gav://org.apache.struts:struts2-osgi-demo-bundle:2.1.8.1",
            "gav://org.apache.struts:struts2-pell-multipart-plugin:2.0.11",
            "gav://org.apache.struts:struts2-config-browser-plugin:2.2.1",
            "gav://org.apache.struts:struts2-dwr-plugin:2.2.1.1",
            "gav://org.apache.struts:struts2-pell-multipart-plugin:2.0.14",
            "gav://org.apache.struts:struts2-jasperreports-plugin:2.0.12",
            "gav://org.apache.struts:struts2-parent:2.2.1.1",
            "gav://org.apache.struts:struts2-jasperreports-plugin:2.0.11",
            "gav://org.apache.struts:struts2-jfreechart-plugin:2.2.1.1",
            "gav://org.apache.struts:struts2-jasperreports-plugin:2.0.14",
            "gav://opensymphony:xwork-tiger-src:1.2.2",
            "gav://org.apache.struts:struts2-testng-plugin:2.1.8.1",
            "gav://opensymphony:xwork-tiger-src:1.2.3",
            "gav://opensymphony:xwork-tiger-src:1.2.1",
            "gav://org.apache.struts:struts2-tiles-plugin:2.1.2",
            "gav://org.apache.struts:struts2-struts1-plugin:2.0.9",
            "gav://org.apache.struts:struts2-tiles-plugin:2.1.8",
            "gav://org.apache.struts:struts2-struts1-plugin:2.0.5",
            "gav://org.apache.struts:struts2-struts1-plugin:2.0.6",
            "gav://org.apache.struts:struts2-struts1-plugin:2.0.8",
            "gav://org.apache.struts:struts2-tiles-plugin:2.1.6",
            "gav://org.apache.struts:struts2-spring-plugin:2.1.8.1",
            "gav://org.apache.struts:struts2-osgi-bundles:2.1.8",
            "gav://org.apache.struts:struts2-tiles-plugin:2.2.1.1",
            "gav://org.apache.struts:struts2-jasperreports-plugin:2.2.1",
            "gav://com.opensymphony:xwork-parent:2.1.6",
            "gav://com.opensymphony:xwork-parent:2.1.5",
            "gav://org.apache.struts:struts2-archetypes:2.2.1",
            "gav://com.opensymphony:xwork-parent:2.1.4",
            "gav://org.apache.struts:struts2-sitegraph-plugin:2.2.1",
            "gav://org.apache.struts:struts2-showcase:2.2.1",
            "gav://org.apache.struts:struts2-mailreader:2.0.14",
            "gav://com.opensymphony:xwork-core:2.1.5",
            "gav://com.opensymphony:xwork-core:2.1.4",
            "gav://org.apache.struts:struts2-osgi-plugin:2.1.8",
            "gav://org.apache.struts:struts2-apps:2.0.11",
            "gav://org.apache.struts:struts2-mailreader:2.1.2",
            "gav://org.apache.struts:struts2-embeddedjsp-plugin:2.2.1.1",
            "gav://org.apache.struts:struts2-showcase:2.2.1.1",
            "gav://org.apache.struts:struts2-spring-plugin:2.0.11",
            "gav://org.apache.struts:struts2-spring-plugin:2.0.12",
            "gav://org.apache.struts:struts2-apps:2.0.12",
            "gav://com.opensymphony:xwork-core:2.1.6",
            "gav://org.apache.struts:struts2-mailreader:2.1.6",
            "gav://org.apache.struts:struts2-apps:2.0.14",
            "gav://com.opensymphony.webwork:com.springsource.com.opensymphony.webwork:2.2.6",
            "gav://org.apache.struts:struts2-mailreader:2.1.8",
            "gav://opensymphony:xwork-tiger-src:1.1.3",
            "gav://opensymphony:xwork-tiger-src:1.1.1",
            "gav://org.apache.struts:struts2-tiles-plugin:2.2.1",
            "gav://org.apache.struts:struts2-struts1-plugin:2.1.8",
            "gav://org.apache.struts:struts2-junit-plugin:2.2.1",
            "gav://org.apache.struts:struts2-spring-plugin:2.0.14",
            "gav://org.apache.struts:struts2-apps:2.1.8.1",
            "gav://org.apache.struts:struts2-spring-plugin:2.0.11.2",
            "gav://org.apache.struts:struts2-struts1-plugin:2.1.6",
            "gav://org.apache.struts:struts2-spring-plugin:2.0.11.1",
            "gav://org.apache.struts:struts2-sitegraph-plugin:2.1.6",
            "gav://org.apache.struts:struts2-struts1-plugin:2.1.2",
            "gav://org.apache.struts:struts2-sitegraph-plugin:2.1.8",
            "gav://org.apache.struts:struts2-struts1-plugin:2.0.14",
            "gav://org.apache.struts:struts2-jasperreports-plugin:2.1.2",
            "gav://org.apache.struts:struts2-jfreechart-plugin:2.0.14",
            "gav://org.apache.struts:struts2-osgi-bundles:2.2.1",
            "gav://org.apache.struts:struts2-mailreader:2.0.11",
            "gav://org.apache.struts:struts2-struts1-plugin:2.0.11",
            "gav://org.apache.struts:struts2-struts1-plugin:2.0.12",
            "gav://org.apache.struts:struts2-mailreader:2.0.12",
            "gav://org.apache.struts:struts2-jfreechart-plugin:2.0.11",
            "gav://org.apache.struts:struts2-sitegraph-plugin:2.1.2",
            "gav://org.apache.struts:struts2-jfreechart-plugin:2.0.12",
            "gav://org.apache.struts:struts2-jasperreports-plugin:2.0.11.2",
            "gav://org.apache.struts:struts2-jasperreports-plugin:2.0.11.1",
            "gav://org.apache.struts:struts2-mailreader:2.0.8",
            "gav://org.apache.struts:struts2-mailreader:2.0.6",
            "gav://opensymphony:xwork-tiger:1.2.3",
            "gav://org.apache.struts:struts2-mailreader:2.0.5",
            "gav://org.apache.struts:struts2-mailreader:2.2.1.1",
            "gav://opensymphony:xwork-tiger:1.2.1",
            "gav://opensymphony:xwork-tiger:1.2.2",
            "gav://org.apache.struts:struts2-mailreader:2.0.9",
            "gav://org.apache.struts:struts2-portlet-plugin:2.1.8.1",
            "gav://org.apache.struts:struts2-pell-multipart-plugin:2.1.8.1",
            "gav://org.apache.struts:struts2-junit-plugin:2.2.1.1",
            "gav://org.apache.struts:struts2-junit-plugin:2.1.2",
            "gav://org.apache.struts:struts2-jsf-plugin:2.2.1",
            "gav://opensymphony:xwork-src:1.1",
            "gav://org.apache.struts:struts2-portlet:2.0.6",
            "gav://opensymphony:xwork-src:1.2",
            "gav://org.apache.struts:struts2-portlet:2.0.5",
            "gav://org.apache.struts:struts2-sitegraph-plugin:2.0.6",
            "gav://org.apache.struts:struts2-spring-plugin:2.0.5",
            "gav://org.apache.struts:struts2-struts1-plugin:2.2.1",
            "gav://org.apache.struts:struts2-junit-plugin:2.1.8",
            "gav://org.apache.struts:struts2-sitegraph-plugin:2.0.8",
            "gav://org.apache.struts:struts2-osgi-plugin:2.2.1",
            "gav://org.apache.struts:struts2-sitegraph-plugin:2.0.9",
            "gav://org.apache.struts:struts2-sitemesh-plugin:2.0.12",
            "gav://org.apache.struts:struts2-junit-plugin:2.1.6",
            "gav://org.apache.struts:struts2-sitemesh-plugin:2.0.11",
            "gav://org.apache.struts:struts2-apps:2.0.11.1",
            "gav://org.apache.struts:struts2-apps:2.0.11.2",
            "gav://org.apache.struts:struts2-sitemesh-plugin:2.0.14",
            "gav://org.apache.struts:struts2-spring-plugin:2.0.8",
            "gav://org.apache.struts:struts2-spring-plugin:2.0.9",
            "gav://org.apache.struts:struts2-spring-plugin:2.0.6",
            "gav://org.apache.struts:struts2-sitegraph-plugin:2.0.5",
            "gav://org.apache.struts:struts2-javatemplates-plugin:2.2.1.1",
            "gav://org.apache.struts:struts2-jasperreports-plugin:2.1.8.1",
            "gav://org.apache.struts:struts2-showcase:2.0.5",
            "gav://org.apache.struts:struts2-showcase:2.0.6",
            "gav://org.apache.struts:struts2-showcase:2.0.8",
            "gav://org.apache.struts:struts2-sitemesh-plugin:2.1.8.1",
            "gav://org.apache.struts:struts2-showcase:2.0.9",
            "gav://org.apache.struts:struts2-jasperreports-plugin:2.1.6",
            "gav://org.apache.struts:struts2-jasperreports-plugin:2.1.8",
            "gav://org.apache.struts:struts2-rest-plugin:2.1.2",
            "gav://org.apache.struts:struts2-rest-plugin:2.2.1.1",
            "gav://org.apache.struts:struts2-portlet:2.0.9",
            "gav://org.apache.struts:struts2-rest-plugin:2.1.6",
            "gav://org.apache.struts:struts2-portlet:2.0.8",
            "gav://org.apache.struts:struts2-rest-plugin:2.1.8",
            "gav://opensymphony:xwork-tiger:1.1.3",
            "gav://opensymphony:xwork-tiger:1.1.1",
            "gav://org.apache.struts:struts2-oval-plugin:2.2.1.1",
            "gav://org.apache.struts:struts2-pell-multipart-plugin:2.0.11.2",
            "gav://org.apache.struts:struts2-tiles-plugin:2.0.8",
            "gav://org.apache.struts:struts2-jsf-plugin:2.1.2",
            "gav://org.apache.struts:struts2-pell-multipart-plugin:2.0.11.1",
            "gav://org.apache.struts:struts2-tiles-plugin:2.0.9",
            "gav://org.apache.struts:struts2-tiles-plugin:2.0.5",
            "gav://org.apache.struts:struts2-tiles-plugin:2.0.6",
            "gav://org.apache.struts:struts2-portlet:2.1.2",
            "gav://org.apache.struts:struts2-spring-plugin:2.1.2",
            "gav://org.apache.struts:struts2-osgi-bundles:2.2.1.1",
            "gav://org.apache.struts:struts2-sitemesh-plugin:2.2.1",
            "gav://org.apache.struts:struts2-spring-plugin:2.1.8",
            "gav://org.apache.struts:struts2-tiles-plugin:2.0.11.1",
            "gav://org.apache.struts:struts2-spring-plugin:2.1.6",
            "gav://org.apache.struts:struts2-showcase:2.1.2",
            "gav://org.apache.struts:struts2-tiles-plugin:2.0.11.2",
            "gav://org.apache.struts:struts2-jfreechart-plugin:2.0.11.1",
            "gav://org.apache.struts:struts2-jfreechart-plugin:2.0.11.2",
            "gav://org.apache.struts:struts2-showcase:2.1.6",
            "gav://org.apache.struts:struts2-parent:2.0.11.2",
            "gav://org.apache.struts:struts2-showcase:2.1.8",
            "gav://org.apache.struts:struts2-parent:2.0.11.1",
            "gav://org.apache.struts:struts2-portlet:2.1.8.1",
            "gav://org.apache.struts:struts2-jasperreports-plugin:2.0.6",
            "gav://org.apache.struts:struts2-jasperreports-plugin:2.0.5",
            "gav://org.apache.struts:struts2-jasperreports-plugin:2.0.8",
            "gav://org.apache.struts:struts2-jasperreports-plugin:2.0.9",
            "gav://org.apache.struts:struts2-rest-plugin:2.2.1",
            "gav://org.apache.archiva.redback:redback-struts2:2.0",
            "gav://org.apache.struts:struts2-portlet:2.1.8",
            "gav://org.apache.struts:struts2-portlet:2.1.6",
            "gav://org.apache.struts:struts2-sitegraph-plugin:2.0.14",
            "gav://org.apache.empire-db:empire-db-struts2:2.1.0-incubating",
            "gav://opensymphony:xwork-src:1.1.3",
            "gav://org.apache.struts:struts2-blank:2.1.2",
            "gav://org.apache.empire-db:empire-db-struts2:2.0.7-incubating",
            "gav://org.apache.struts:struts2-sitemesh-plugin:2.1.6",
            "gav://org.apache.struts:struts2-blank:2.1.8",
            "gav://org.apache.struts:struts2-sitemesh-plugin:2.1.8",
            "gav://org.apache.struts:struts2-blank:2.1.6",
            "gav://opensymphony:xwork-src:1.1.1",
            "gav://org.apache.struts:struts2-archetype-blank:2.2.1",
            "gav://org.apache.struts:struts2-pell-multipart-plugin:2.0.8",
            "gav://org.apache.struts:struts2-archetype-starter:2.2.1",
            "gav://org.apache.struts:struts2-pell-multipart-plugin:2.0.6",
            "gav://org.apache.struts:struts2-pell-multipart-plugin:2.0.5",
            "gav://org.apache.struts:struts2-sitemesh-plugin:2.1.2",
            "gav://org.apache.struts:struts2-pell-multipart-plugin:2.0.9",
            "gav://org.apache.struts:struts2-assembly:2.2.1",
            "gav://org.apache.struts:struts2-codebehind-plugin:2.0.14",
            "gav://org.apache.struts:struts2-convention-plugin:2.1.8",
            "gav://org.apache.struts:struts2-jboss-blank:2.2.1",
            "gav://org.apache.struts:struts2-codebehind-plugin:2.0.12",
            "gav://org.apache.struts:struts2-convention-plugin:2.1.6",
            "gav://org.apache.struts:struts2-codebehind-plugin:2.0.11",
            "gav://org.apache.struts:struts2-parent:2.1.8.1",
            "gav://org.apache.struts:struts2-core:2.2.1",
            "gav://org.apache.struts:struts2-codebehind-plugin:2.1.6",
            "gav://org.apache.struts:struts2-codebehind-plugin:2.1.2",
            "gav://org.apache.struts:struts2-apps:2.1.2",
            "gav://org.apache.struts:struts2-json-plugin:2.1.8",
            "gav://org.apache.struts:struts2-archetype-portlet:2.2.1.1",
            "gav://org.apache.struts:struts2-apps:2.1.6",
            "gav://org.apache.struts:struts2-portlet-plugin:2.2.1",
            "gav://org.apache.struts:struts2-core:2.1.8.1",
            "gav://org.apache.struts:struts2-apps:2.1.8",
            "gav://opensymphony:xwork-src:1.2.2",
            "gav://opensymphony:xwork-src:1.2.3",
            "gav://org.apache.struts:struts2-sitemesh-plugin:2.0.5",
            "gav://org.apache.struts:struts2-sitemesh-plugin:2.0.6",
            "gav://org.apache.struts:struts2-sitemesh-plugin:2.0.8",
            "gav://org.apache.struts:struts2-blank:2.0.9",
            "gav://org.apache.struts:struts2-blank:2.0.8",
            "gav://org.apache.struts:struts2-sitemesh-plugin:2.0.9",
            "gav://org.apache.struts:struts2-blank:2.0.6",
            "gav://opensymphony:xwork-src:1.2.1",
            "gav://org.apache.struts:struts2-blank:2.0.5",
            "gav://org.apache.struts:struts2-junit-plugin:2.1.8.1",
            "gav://org.apache.struts:struts2-core:2.0.11.1",
            "gav://org.apache.struts:struts2-dojo-plugin:2.2.1",
            "gav://org.apache.struts:struts2-json-plugin:2.2.1",
            "gav://org.apache.struts:struts2-core:2.0.11.2",
            "gav://org.apache.struts:struts2-osgi-admin-bundle:2.2.1",
            "gav://org.apache.struts:struts2-plugins:2.2.1.1",
            "gav://org.apache.struts:struts2-codebehind-plugin:2.2.1",
            "gav://org.apache.struts:struts2-apps:2.0.8",
            "gav://org.apache.struts:struts2-portlet-plugin:2.1.6",
            "gav://org.apache.struts:struts2-apps:2.0.6",
            "gav://org.apache.struts:struts2-apps:2.0.5",
            "gav://org.apache.struts:struts2-portlet-plugin:2.1.2",
            "gav://org.apache.struts:struts2-osgi-plugin:2.1.8.1",
            "gav://org.apache.struts:struts2-apps:2.0.9",
            "gav://org.apache.struts:struts2-showcase:2.1.8.1",
            "gav://org.apache.struts:struts2-portlet:2.0.11.2",
            "gav://org.apache.struts:struts2-portlet:2.0.11.1",
            "gav://org.apache.struts:struts2-portlet-plugin:2.1.8",
            "gav://org.apache.struts:struts2-osgi-bundles:2.1.8.1",
            "gav://org.apache.struts:struts2-codebehind-plugin:2.1.8.1",
            "gav://org.apache.struts:struts2-json-plugin:2.2.1.1",
            "gav://org.apache.struts:struts2-blank:2.0.11.2",
            "gav://org.apache.struts:struts2-pell-multipart-plugin:2.2.1",
            "gav://org.apache.struts:struts2-blank:2.0.11.1",
            "gav://com.opensymphony:xwork-assembly:2.1.4",
            "gav://org.apache.struts:struts2-plexus-plugin:2.0.11.2",
            "gav://com.opensymphony:xwork-assembly:2.1.6",
            "gav://org.apache.struts:struts2-javatemplates-plugin:2.2.1",
            "gav://org.apache.struts:struts2-dojo-plugin:2.1.2",
            "gav://org.apache.struts:struts2-convention-plugin:2.1.8.1",
            "gav://com.opensymphony:xwork-assembly:2.1.5",
            "gav://org.apache.struts:struts2-plexus-plugin:2.0.11.1",
            "gav://org.apache.struts:struts2-dojo-plugin:2.1.6",
            "gav://org.apache.struts:struts2-sitegraph-plugin:2.1.8.1",
            "gav://org.apache.struts:struts2-codebehind-plugin:2.0.8",
            "gav://org.apache.struts:struts2-codebehind-plugin:2.0.9",
            "gav://org.apache.struts:struts2-dojo-plugin:2.1.8",
            "gav://org.apache.struts:struts2-core:2.0.9",
            "gav://org.apache.struts:struts2-core:2.0.8",
            "gav://org.apache.struts:struts2-core:2.0.5",
            "gav://org.apache.struts:struts2-core:2.0.6",
            "gav://org.apache.struts:struts2-archetype-dbportlet:2.2.1.1",
            "gav://org.apache.struts:struts2-struts1-plugin:2.2.1.1",
            "gav://opensymphony:xwork-tiger-src:2.0-beta-1",
            "gav://org.apache.struts:struts2-pell-multipart-plugin:2.2.1.1",
            "gav://org.apache.struts:struts2-assembly:2.0.14",
            "gav://org.apache.archiva.redback:redback-struts2-content:2.0",
            "gav://org.apache.struts:struts2-assembly:2.0.11",
            "gav://org.apache.struts:struts2-blank:2.2.1",
            "gav://org.apache.struts:struts2-pell-multipart-plugin:2.1.2",
            "gav://org.apache.struts:struts2-blank:2.1.8.1",
            "gav://org.apache.struts:struts2-pell-multipart-plugin:2.1.6",
            "gav://org.apache.struts:struts2-pell-multipart-plugin:2.1.8",
            "gav://org.apache.struts:struts2-codebehind-plugin:2.1.8",
            "gav://org.apache.struts:struts2-convention-plugin:2.2.1",
            "gav://org.apache.struts:struts2-core:2.1.8",
            "gav://org.apache.struts:struts2-javatemplates-plugin:2.1.8",
            "gav://org.apache.struts:struts2-portlet-plugin:2.2.1.1",
            "gav://org.apache.struts:struts2-javatemplates-plugin:2.1.6",
            "gav://org.apache.struts:struts2-embeddedjsp-plugin:2.1.8.1",
            "gav://org.apache.struts:struts2-core:2.1.6",
            "gav://org.apache.struts:struts2-sitemesh-plugin:2.0.11.2",
            "gav://org.apache.struts:struts2-sitemesh-plugin:2.0.11.1",
            "gav://org.apache.struts:struts2-core:2.1.2",
            "gav://org.apache.struts:struts2-config-browser-plugin:2.0.11.2",
            "gav://org.apache.struts:struts2-codebehind-plugin:2.0.6",
            "gav://org.apache.struts:struts2-config-browser-plugin:2.0.11.1",
            "gav://org.apache.struts:struts2-codebehind-plugin:2.0.5",
            "gav://org.apache.struts:struts2-plugins:2.0.12",
            "gav://org.apache.struts:struts2-plugins:2.0.11",
            "gav://org.apache.struts:struts2-mailreader:2.2.1",
            "gav://org.apache.struts:struts2-plugins:2.0.14",
            "gav://com.opensymphony:xwork-plugins:2.1.4",
            "gav://org.apache.struts:struts2-osgi-demo-bundle:2.1.8",
            "gav://com.opensymphony:xwork-plugins:2.1.5",
            "gav://com.opensymphony:xwork-plugins:2.1.6",
            "gav://org.apache.struts:struts2-archetype-starter:2.0.11.2",
            "gav://org.apache.struts:struts2-blank:2.2.1.1",
            "gav://org.apache.struts:struts2-plexus-plugin:2.0.5",
            "gav://org.apache.struts:struts2-embeddedjsp-plugin:2.2.1",
            "gav://org.apache.struts:struts2-plexus-plugin:2.0.6",
            "gav://org.apache.struts:struts2-jsf-plugin:2.0.14",
            "gav://org.apache.struts:struts2-jsf-plugin:2.0.11",
            "gav://org.apache.struts:struts2-plexus-plugin:2.0.8",
            "gav://org.apache.struts:struts2-jsf-plugin:2.0.12",
            "gav://org.apache.struts:struts2-plexus-plugin:2.0.9",
            "gav://org.apache.struts:struts2-dwr-plugin:2.1.8.1",
            "gav://org.apache.struts:struts2-tiles-plugin:2.0.12",
            "gav://org.apache.struts:struts2-tiles-plugin:2.0.11",
            "gav://org.apache.struts:struts2-tiles-plugin:2.0.14",
            "gav://org.apache.struts:struts2-plugins:2.1.8.1",
            "gav://org.apache.struts:struts2-parent:2.1.2",
            "gav://org.apache.struts:struts2-parent:2.1.8",
            "gav://org.apache.struts:struts2-oval-plugin:2.1.8.1",
            "gav://org.apache.struts:struts2-parent:2.1.6",
            "gav://org.apache.struts:struts2-osgi-demo-bundle:2.2.1",
            "gav://com.opensymphony.xwork:com.springsource.com.opensymphony.xwork:1.2.3",
            "gav://org.apache.struts:struts2-core:2.2.1.1",
            "gav://org.apache.struts:struts2-plexus-plugin:2.1.2",
            "gav://opensymphony:webwork-src:2.2",
            "gav://org.apache.struts:struts2-plexus-plugin:2.1.6",
            "gav://org.apache.struts:struts2-plexus-plugin:2.1.8",
            "gav://org.apache.struts:struts2-sitegraph-plugin:2.2.1.1",
            "gav://org.apache.struts:struts2-core:2.0.12",
            "gav://org.apache.struts:struts2-core:2.0.11",
            "gav://org.apache.struts:struts2-core:2.0.14",
            "gav://org.apache.struts:struts2-assembly:2.2.1.1",
            "gav://org.apache.struts:struts2-mailreader:2.0.11.1",
            "gav://org.apache.struts:struts2-plugins:2.0.11.1",
            "gav://org.apache.struts:struts2-mailreader:2.0.11.2",
            "gav://org.apache.struts:struts2-plugins:2.0.11.2",
            "gav://org.apache.struts:struts2-parent:2.0.9",
            "gav://org.apache.struts:struts2-api:2.0.5",
            "gav://opensymphony:xwork-src:2.0-beta-1",
            "gav://org.apache.struts:struts2-parent:2.0.6",
            "gav://org.apache.struts:struts2-archetype-plugin:2.2.1.1",
            "gav://opensymphony:xwork-src:2.0-beta-2",
            "gav://org.apache.struts:struts2-parent:2.0.5",
            "gav://org.apache.struts:struts2-parent:2.0.8",
            "gav://opensymphony:webwork-src:2.2.3",
            "gav://org.apache.struts:struts2-plugins:2.0.9",
            "gav://opensymphony:webwork-src:2.2.4",
            "gav://opensymphony:webwork-src:2.2.1",
            "gav://opensymphony:webwork-src:2.2.2",
            "gav://org.apache.struts:struts2-struts1-plugin:2.0.11.1",
            "gav://org.apache.struts:struts2-struts1-plugin:2.0.11.2",
            "gav://opensymphony:xwork-tiger:1.2",
            "gav://opensymphony:xwork-tiger:1.1",
            "gav://org.apache.struts:struts2-plugins:2.0.6",
            "gav://org.apache.struts:struts2-plugins:2.0.5",
            "gav://org.apache.struts:struts2-plugins:2.0.8",
            "gav://opensymphony:webwork-src:2.2.5",
            "gav://org.apache.struts:struts2-sitemesh-plugin:2.2.1.1",
            "gav://org.apache.struts:struts2-struts1-plugin:2.1.8.1",
            "gav://com.opensymphony.xwork:com.springsource.com.opensymphony.xwork:1.1.1",
            "gav://org.apache.struts:struts2-archetype-dbportlet:2.2.1",
            "gav://opensymphony:xwork-tiger-src:1.2",
            "gav://org.apache.struts:struts2-plexus-plugin:2.2.1",
            "gav://opensymphony:xwork-tiger-src:1.1",
            "gav://org.apache.struts:struts2-plexus-plugin:2.0.14",
            "gav://org.apache.struts:struts2-rest-showcase:2.2.1",
            "gav://org.apache.struts:struts2-blank:2.0.12",
            "gav://org.apache.archiva.redback:redback-struts2-integration:2.0",
            "gav://org.apache.struts:struts2-javatemplates-plugin:2.1.8.1",
            "gav://org.apache.struts:struts2-blank:2.0.11",
            "gav://org.apache.struts:struts2-showcase:2.0.11.1",
            "gav://org.apache.struts:struts2-osgi-admin-bundle:2.1.8",
            "gav://org.apache.struts:struts2-showcase:2.0.11.2",
            "gav://org.apache.struts:struts2-plugins:2.1.8",
            "gav://org.apache.struts:struts2-dwr-plugin:2.1.2",
            "gav://org.apache.struts:struts2-rest-showcase:2.1.8.1",
            "gav://org.apache.struts.xwork:xwork-core:2.2.1.1",
            "gav://org.apache.struts:struts2-plugins:2.1.2",
            "gav://org.apache.struts:struts2-codebehind-plugin:2.2.1.1",
            "gav://org.apache.struts:struts2-plugins:2.1.6",
            "gav://org.apache.struts:struts2-jsf-plugin:2.1.8.1",
            "gav://org.apache.struts:struts2-portlet:2.0.12",
            "gav://org.apache.struts:struts2-embeddedjsp-plugin:2.1.8",
            "gav://org.apache.struts:struts2-portlet:2.0.14",
            "gav://org.apache.struts:struts2-archetype-plugin:2.2.1",
            "gav://org.apache.struts:struts2-jfreechart-plugin:2.1.8.1",
            "gav://org.apache.struts:struts2-config-browser-plugin:2.1.8.1",
            "gav://org.apache.struts:struts2-osgi-plugin:2.2.1.1",
            "gav://org.apache.struts:struts2-portlet:2.0.11",
            "gav://org.apache.struts:struts2-convention-plugin:2.2.1.1",
            "gav://org.apache.struts:struts2-mailreader:2.1.8.1",
            "gav://org.apache.struts:struts2-parent:2.0.12",
            "gav://org.apache.struts:struts2-testng-plugin:2.2.1",
            "gav://org.apache.struts:struts2-parent:2.0.11",
            "gav://org.apache.struts:struts2-rest-showcase:2.1.2",
            "gav://org.apache.struts:struts2-parent:2.0.14",
            "gav://org.apache.struts:struts2-rest-showcase:2.1.6",
            "gav://org.apache.struts:struts2-rest-showcase:2.1.8",
            "gav://org.apache.struts:struts2-sitegraph-plugin:2.0.11.1",
            "gav://org.apache.struts:struts2-sitegraph-plugin:2.0.11.2",
            "gav://org.apache.struts:struts2-dwr-plugin:2.1.8",
            "gav://org.apache.struts:struts2-dwr-plugin:2.1.6",
            "gav://org.apache.struts:struts2-archetypes:2.2.1.1",
            "gav://org.apache.empire-db:empire-db-struts2:2.0.5-incubating",
            "gav://org.apache.struts:struts2-dojo-plugin:2.2.1.1",
            "gav://org.apache.struts:struts2-plugins:2.2.1",
            "gav://org.apache.struts:struts2-json-plugin:2.1.8.1",
            "gav://org.apache.struts:struts2-apps:2.2.1",
            "gav://org.apache.struts:struts2-blank:2.0.14",
            "gav://org.apache.struts.xwork:xwork-core:2.2.1",
            "gav://org.apache.struts:struts2-plexus-plugin:2.2.1.1",
            "gav://org.apache.struts:struts2-jfreechart-plugin:2.0.9",
            "gav://org.apache.struts:struts2-osgi-admin-bundle:2.2.1.1",
            "gav://org.apache.struts:struts2-plexus-plugin:2.0.11",
            "gav://org.apache.struts:struts2-plexus-plugin:2.0.12",
            "gav://org.apache.struts:struts2-showcase:2.0.14",
            "gav://org.apache.struts:struts2-apps:2.2.1.1",
            "gav://org.apache.struts:struts2-showcase:2.0.11",
            "gav://org.apache.struts:struts2-archetype-portlet:2.2.1",
            "gav://org.apache.struts:struts2-showcase:2.0.12"
          ],
          "impacted_artifacts": [
            {
              "name": "pipeline-xray",
              "path": "LOCAL/builds/",
              "sha256": "f8fa924e200c8412075de89d867f1b8dd5aa343938754fb652f84ff472a75b97",
              "sha1": "",
              "depth": 0,
              "display_name": "pipeline-xray:294",
              "pkg_type": "Build",
              "parent_sha": "f8fa924e200c8412075de89d867f1b8dd5aa343938754fb652f84ff472a75b97",
              "impact_path": "",
              "infected_file": {
                "name": "xwork-2.0.6.jar",
                "path": "LOCAL/libs-release-local/com/opensymphony/xwork/2.0.6/",
                "sha256": "4baa7d0a7203ceeefa063a3782de733d417c4eae5ee5c0ab06e9e9e4607c004a",
                "sha1": "",
                "depth": 0,
                "parent_sha": "4baa7d0a7203ceeefa063a3782de733d417c4eae5ee5c0ab06e9e9e4607c004a",
                "display_name": "com.opensymphony:xwork:2.0.6",
                "pkg_type": "Maven"
              }
            }
          ]
        }
      ],
      "top_severity": "Major",
      "watch_name": "all"
    },
    {
      "created": "2017-01-11T12:31:35.816224803+02:00",
      "issues": [
        {
          "severity": "Major",
          "type": "security",
          "provider": "JFrog",
          "created": "2011-05-13T16:05:45+03:00",
          "summary": "CWE-200 Information Exposure",
          "description": "XWork 2.2.1 in Apache Struts 2.2.1, and OpenSymphony XWork in OpenSymphony WebWork, allows remote attackers to obtain potentially sensitive information about internal Java class paths via vectors involving an s:submit element and a nonexistent method, a different vulnerability than CVE-2011-1772.3.",
          "cve": "CVE-2011-2088",
          "component_ids": [
            "gav://opensymphony:xwork:1.0.3.6",
            "gav://opensymphony:xwork:1.0.3.5",
            "gav://opensymphony:xwork:1.0.3.4",
            "gav://opensymphony:xwork:1.0.3.3",
            "gav://opensymphony:xwork:1.1.1",
            "gav://opensymphony:xwork:1.0.3.2",
            "gav://opensymphony:xwork:1.0.3.1",
            "gav://opensymphony:xwork:1.1.3",
            "gav://opensymphony:webwork:2.1.5-atlassian-2",
            "gav://opensymphony:webwork:2.1.5-atlassian-3",
            "gav://opensymphony:webwork:2.2.1",
            "gav://opensymphony:webwork:2.2.2",
            "gav://opensymphony:webwork:2.2.3",
            "gav://opensymphony:xwork:2.0-beta-3-j4",
            "gav://opensymphony:xwork:2.0.1-j4",
            "gav://com.opensymphony:webwork:2.2.6-atlassian-3",
            "gav://opensymphony:webwork:2.2.4",
            "gav://opensymphony:webwork:2.2.5",
            "gav://com.opensymphony:webwork:2.2.6-atlassian-2",
            "gav://opensymphony:webwork:2.2.7",
            "gav://opensymphony:xwork:2.0.3",
            "gav://opensymphony:xwork:2.0.2",
            "gav://com.opensymphony:xwork:2.0.7",
            "gav://com.opensymphony:xwork:2.0.5",
            "gav://com.opensymphony:xwork:2.0.6",
            "gav://opensymphony:webwork:2.2.7-atlassian-25",
            "gav://opensymphony:webwork:2.2.7-atlassian-27",
            "gav://opensymphony:xwork:2.0.0-j4",
            "gav://opensymphony:xwork:1.2.2",
            "gav://opensymphony:xwork:1.2.1",
            "gav://opensymphony:webwork:2.2.7-atlassian-23",
            "gav://opensymphony:xwork:1.2.3",
            "gav://opensymphony:webwork:2.2.7-atlassian-24",
            "gav://opensymphony:xwork:2.0-beta-2",
            "gav://opensymphony:xwork:2.0-RC1",
            "gav://opensymphony:xwork:2.0-beta-3",
            "gav://opensymphony:xwork:2.0-RC1-j4",
            "gav://opensymphony:xwork:2.0-beta-1",
            "gav://opensymphony:xwork:2.0.1",
            "gav://opensymphony:xwork:2.0.0",
            "gav://opensymphony:xwork:1.1",
            "gav://opensymphony:xwork:1.2",
            "gav://opensymphony:xwork:1.0",
            "gav://opensymphony:webwork:2.2.7-atlassian-29",
            "gav://opensymphony:webwork:1.4-atlassian-26-testmvn2-1",
            "gav://opensymphony:webwork:2.2.7-atlassian-1",
            "gav://opensymphony:webwork:2.2.7-atlassian-2",
            "gav://opensymphony:webwork:2.2.7-atlassian-3",
            "gav://opensymphony:webwork:2.2.7-atlassian-5",
            "gav://opensymphony:webwork:2.2.7-atlassian-6",
            "gav://com.opensymphony:webwork:2.2.6",
            "gav://opensymphony:webwork:1.4-atlassian-9",
            "gav://opensymphony:webwork:2.1",
            "gav://opensymphony:webwork:1.4-atlassian-4",
            "gav://opensymphony:webwork:2.0",
            "gav://opensymphony:webwork:1.4-atlassian-2",
            "gav://opensymphony:webwork:2.2",
            "gav://opensymphony:webwork:1.4-atlassian-1",
            "gav://com.opensymphony:xwork:2.0.4",
            "gav://opensymphony:xwork:1.2.5-rc1",
            "gav://opensymphony:xwork:1.0.1",
            "gav://opensymphony:webwork:1.4-atlassian-19",
            "gav://opensymphony:xwork:1.2.5-atlassian-1",
            "gav://opensymphony:xwork:1.2.5-atlassian-2",
            "gav://opensymphony:xwork:1.2.5-atlassian-4",
            "gav://opensymphony:xwork:1.0.4",
            "gav://opensymphony:xwork:1.0.3",
            "gav://opensymphony:xwork:1.2.5-atlassian-5",
            "gav://opensymphony:xwork:1.2.5-atlassian-6",
            "gav://opensymphony:xwork:1.2.5-atlassian-7",
            "gav://opensymphony:xwork:1.0.5",
            "gav://opensymphony:webwork:1.4-atlassian-10",
            "gav://opensymphony:webwork:1.4-atlassian-11",
            "gav://opensymphony:webwork:1.4-atlassian-12",
            "gav://opensymphony:webwork:1.4-atlassian-13",
            "gav://opensymphony:webwork:1.4-atlassian-15",
            "gav://opensymphony:webwork:1.4-atlassian-16",
            "gav://opensymphony:webwork:1.4-atlassian-17",
            "gav://opensymphony:xwork:1.2.3-20070717",
            "gav://opensymphony:webwork:1.4-atlassian-30",
            "gav://opensymphony:webwork:12Dec05-jiratld",
            "gav://opensymphony:xwork:1.2.5-atlassian-8",
            "gav://opensymphony:webwork:1.4",
            "gav://opensymphony:webwork:2.1.3",
            "gav://opensymphony:webwork:2.1.4",
            "gav://com.opensymphony:xwork:2.1.2",
            "gav://opensymphony:webwork:2.1.5",
            "gav://com.opensymphony:xwork:2.1.3",
            "gav://opensymphony:webwork:2.1.6",
            "gav://com.opensymphony:xwork:2.1.0",
            "gav://opensymphony:webwork:2.1.7",
            "gav://com.opensymphony:xwork:2.1.1",
            "gav://opensymphony:webwork:1.4-atlassian-22",
            "gav://opensymphony:webwork:1.4-atlassian-23",
            "gav://opensymphony:webwork:1.4-atlassian-24",
            "gav://opensymphony:webwork:1.4-atlassian-25",
            "gav://opensymphony:webwork:1.4-atlassian-27",
            "gav://opensymphony:xwork-src:1.1.3",
            "gav://opensymphony:xwork-tiger-src:1.2.2",
            "gav://opensymphony:xwork-tiger-src:1.2.3",
            "gav://opensymphony:xwork-tiger-src:1.2.1",
            "gav://org.apache.struts:struts2-embeddedjsp-plugin:2.2.1",
            "gav://opensymphony:xwork-src:1.1.1",
            "gav://org.apache.struts:struts2-archetype-blank:2.2.1",
            "gav://org.apache.struts:struts2-portlet:2.2.1",
            "gav://org.apache.struts:struts2-archetype-starter:2.2.1",
            "gav://org.apache.struts:struts2-spring-plugin:2.2.1",
            "gav://org.apache.struts:struts2-jasperreports-plugin:2.2.1",
            "gav://com.opensymphony:xwork-parent:2.1.6",
            "gav://com.opensymphony:xwork-parent:2.1.5",
            "gav://org.apache.struts:struts2-archetypes:2.2.1",
            "gav://org.apache.struts:struts2-assembly:2.2.1",
            "gav://com.opensymphony:xwork-parent:2.1.4",
            "gav://org.apache.struts:struts2-jboss-blank:2.2.1",
            "gav://org.apache.struts:struts2-sitegraph-plugin:2.2.1",
            "gav://org.apache.struts:struts2-showcase:2.2.1",
            "gav://com.opensymphony:xwork-core:2.1.5",
            "gav://org.apache.struts:struts2-core:2.2.1",
            "gav://com.opensymphony:xwork-core:2.1.4",
            "gav://org.apache.struts:struts2-portlet-plugin:2.2.1",
            "gav://com.opensymphony:xwork-core:2.1.6",
            "gav://org.apache.struts:struts2-osgi-demo-bundle:2.2.1",
            "gav://com.opensymphony.webwork:com.springsource.com.opensymphony.webwork:2.2.6",
            "gav://com.opensymphony.xwork:com.springsource.com.opensymphony.xwork:1.2.3",
            "gav://opensymphony:xwork-tiger-src:1.1.3",
            "gav://opensymphony:xwork-src:1.2.2",
            "gav://opensymphony:xwork-src:1.2.3",
            "gav://opensymphony:xwork-tiger-src:1.1.1",
            "gav://org.apache.struts:struts2-tiles-plugin:2.2.1",
            "gav://org.apache.struts:struts2-junit-plugin:2.2.1",
            "gav://opensymphony:webwork-src:2.2",
            "gav://opensymphony:xwork-src:1.2.1",
            "gav://org.apache.struts:struts2-dojo-plugin:2.2.1",
            "gav://org.apache.struts:struts2-json-plugin:2.2.1",
            "gav://org.apache.struts:struts2-osgi-bundles:2.2.1",
            "gav://org.apache.struts:struts2-osgi-admin-bundle:2.2.1",
            "gav://opensymphony:xwork-src:2.0-beta-1",
            "gav://opensymphony:xwork-src:2.0-beta-2",
            "gav://opensymphony:webwork-src:2.2.3",
            "gav://opensymphony:webwork-src:2.2.4",
            "gav://opensymphony:webwork-src:2.2.1",
            "gav://org.apache.struts:struts2-dwr-plugin:2.2.1",
            "gav://opensymphony:webwork-src:2.2.2",
            "gav://opensymphony:xwork-tiger:1.2",
            "gav://org.apache.struts:struts2-oval-plugin:2.2.1",
            "gav://opensymphony:xwork-tiger:1.1",
            "gav://org.apache.struts:struts2-codebehind-plugin:2.2.1",
            "gav://opensymphony:webwork-src:2.2.5",
            "gav://opensymphony:xwork-tiger:1.2.3",
            "gav://opensymphony:xwork-tiger:1.2.1",
            "gav://com.opensymphony.xwork:com.springsource.com.opensymphony.xwork:1.1.1",
            "gav://opensymphony:xwork-tiger:1.2.2",
            "gav://org.apache.struts:struts2-archetype-dbportlet:2.2.1",
            "gav://opensymphony:xwork-tiger-src:1.2",
            "gav://org.apache.struts:struts2-plexus-plugin:2.2.1",
            "gav://opensymphony:xwork-tiger-src:1.1",
            "gav://org.apache.struts:struts2-jfreechart-plugin:2.2.1",
            "gav://org.apache.struts:struts2-jsf-plugin:2.2.1",
            "gav://opensymphony:xwork-src:1.1",
            "gav://org.apache.struts:struts2-pell-multipart-plugin:2.2.1",
            "gav://opensymphony:xwork-src:1.2",
            "gav://com.opensymphony:xwork-assembly:2.1.4",
            "gav://org.apache.struts:struts2-rest-showcase:2.2.1",
            "gav://com.opensymphony:xwork-assembly:2.1.6",
            "gav://org.apache.struts:struts2-struts1-plugin:2.2.1",
            "gav://org.apache.struts:struts2-javatemplates-plugin:2.2.1",
            "gav://com.opensymphony:xwork-assembly:2.1.5",
            "gav://org.apache.struts:struts2-osgi-plugin:2.2.1",
            "gav://org.apache.struts:struts2-archetype-convention:2.2.1",
            "gav://opensymphony:xwork-tiger:2.0-beta-1",
            "gav://opensymphony:xwork-tiger-src:2.0-beta-1",
            "gav://org.apache.struts:struts2-archetype-plugin:2.2.1",
            "gav://opensymphony:xwork-tiger:1.1.3",
            "gav://opensymphony:xwork-tiger:1.1.1",
            "gav://org.apache.struts:struts2-blank:2.2.1",
            "gav://org.apache.struts:struts2-testng-plugin:2.2.1",
            "gav://org.apache.struts:struts2-parent:2.2.1",
            "gav://org.apache.struts:struts2-sitemesh-plugin:2.2.1",
            "gav://org.apache.struts:struts2-config-browser-plugin:2.2.1",
            "gav://org.apache.struts:struts2-convention-plugin:2.2.1",
            "gav://org.apache.struts:struts2-plugins:2.2.1",
            "gav://org.apache.struts:struts2-rest-plugin:2.2.1",
            "gav://org.apache.struts:struts2-apps:2.2.1",
            "gav://org.apache.struts.xwork:xwork-core:2.2.1",
            "gav://org.apache.struts:struts2-mailreader:2.2.1",
            "gav://com.opensymphony:xwork-plugins:2.1.4",
            "gav://com.opensymphony:xwork-plugins:2.1.5",
            "gav://com.opensymphony:xwork-plugins:2.1.6",
            "gav://org.apache.struts:struts2-archetype-portlet:2.2.1"
          ],
          "impacted_artifacts": [
            {
              "name": "pipeline-xray",
              "path": "LOCAL/builds/",
              "sha256": "f8fa924e200c8412075de89d867f1b8dd5aa343938754fb652f84ff472a75b97",
              "sha1": "",
              "depth": 0,
              "display_name": "pipeline-xray:294",
              "pkg_type": "Build",
              "parent_sha": "f8fa924e200c8412075de89d867f1b8dd5aa343938754fb652f84ff472a75b97",
              "impact_path": "",
              "infected_file": {
                "name": "xwork-2.0.6.jar",
                "path": "LOCAL/libs-release-local/com/opensymphony/xwork/2.0.6/",
                "sha256": "4baa7d0a7203ceeefa063a3782de733d417c4eae5ee5c0ab06e9e9e4607c004a",
                "sha1": "",
                "depth": 0,
                "parent_sha": "4baa7d0a7203ceeefa063a3782de733d417c4eae5ee5c0ab06e9e9e4607c004a",
                "display_name": "com.opensymphony:xwork:2.0.6",
                "pkg_type": "Maven"
              }
            }
          ]
        }
      ],
      "top_severity": "Major",
      "watch_name": "my-jenkins"
    }
  ],
  "licenses": [
    {
      "name": "Unknown",
      "components": [
        "build://pipeline-xray:294",
        "gav://com.opensymphony:xwork:2.0.6"
      ],
      "full_name": "Unknown license",
      "more_info_url": null
    },
    {
      "name": "MIT",
      "components": [
        "gav://org.webjars.bower:vf-angular-ui-router:1.0.0-beta.3"
      ],
      "full_name": "MIT license ",
      "more_info_url": [
        "https://opensource.org/licenses/MIT"
      ]
    }
  ]
}
`

const VulnerabilityXrayReportRequestResponse = `
{
  "report_id": 777,
  "status": "pending"
}
`

const LicensesXrayReportRequestResponse = `
{
  "report_id": 888,
  "status": "pending"
}
`

const VulnerabilityReportStatusResponse = `
{
  "id": 301,
  "name": "test-generic",
  "report_type": "vulnerability",
  "status": "completed",
  "total_artifacts": 4,
  "num_of_processed_artifacts": 4,
  "progress": 100,
  "number_of_rows": 64,
  "start_time": "2021-09-03T21:17:41Z",
  "end_time": "2021-09-03T21:17:42Z",
  "author": "test"
}
`

const LicensesReportStatusResponse = `
{
  "id": 301,
  "name": "test-generic",
  "report_type": "license",
  "status": "completed",
  "total_artifacts": 4,
  "num_of_processed_artifacts": 4,
  "progress": 100,
  "number_of_rows": 64,
  "start_time": "2021-09-03T21:17:41Z",
  "end_time": "2021-09-03T21:17:42Z",
  "author": "test"
}
`

const XrayReportDeleteResponse = `
{
  "info": "report deleted successfully"
}
`

const VulnerabilityReportDetailsResponse = `
{
  "total_rows": 70,
  "rows": [
    {
      "cves": [
        {
          "cve": "CVE-2021-37136"
        },
        {
          "cvss_v2_score": 7.1,
          "cvss_v2_vector": "CVSS:2.0/AV:N/AC:M/Au:N/C:N/I:N/A:C",
          "cvss_v3_score": 7.5,
          "cvss_v3_vector": "CVSS:3.0/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H"
        }
      ],
      "cvss2_max_score": 7.1,
      "cvss3_max_score": 7.5,
      "summary": "Netty Bzip2Decoder Class Missing Decompressed Data Allocation Limitation Memory Exhaustion Remote DoS",
      "severity": "High",
      "severity_source": "CVSS V3 from RBS",
      "vulnerable_component": "gav://io.netty:netty-codec:4.1.66.Final",
      "impacted_artifact": "docker://elasticsearch/7.15.0-test2:latest",
      "impact_path": [
        "docker://elasticsearch/7.15.0-test2:latest",
        "generic://sha256:a53372ba228046f81171efd357179b7b02a95acccee17925e3d0295829cb42ea/sha256__a53372ba228046f81171efd357179b7b02a95acccee17925e3d0295829cb42ea.tar.gz",
        "generic://sha256:e1a7a6f8262f89cad679f309ec7875c9a995099ee0fee1a1380ee71692657a4c/elasticsearch-sql-cli-7.15.0.jar",
        "gav://io.netty:netty-codec:4.1.66.Final"
      ],
      "path": "forks-release-local/elasticsearch/7.15.0-test2/latest/",
      "fixed_versions": [
        "4.1.68.Final"
      ],
      "published": "2021-09-12T14:10:55Z",
      "issue_id": "XRAY-184735",
      "package_type": "maven",
      "provider": "JFrog",
      "description": "Netty contains a flaw in the Bzip2Decoder class that is triggered as memory allocations for decompressed data can't be properly limited. This may allow a remote attacker to exhaust available memory resources and cause a denial of service.",
      "references": [
        "https://github.com/netty/netty/commit/41d3d61a61608f2223bb364955ab2045dd5e4020",
        "https://github.com/netty/netty/security/advisories/GHSA-grg4-wf29-r9vv",
        "http://cve.mitre.org/cgi-bin/cvename.cgi?name=2021-37136"
      ]
    },
    {
      "cves": [
        {
          "cve": "CVE-2021-37136"
        },
        {
          "cvss_v2_score": 7.1,
          "cvss_v2_vector": "CVSS:2.0/AV:N/AC:M/Au:N/C:N/I:N/A:C",
          "cvss_v3_score": 7.5,
          "cvss_v3_vector": "CVSS:3.0/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H"
        }
      ],
      "cvss2_max_score": 7.1,
      "cvss3_max_score": 7.5,
      "summary": "Netty Bzip2Decoder Class Missing Decompressed Data Allocation Limitation Memory Exhaustion Remote DoS",
      "severity": "High",
      "severity_source": "CVSS V3 from RBS",
      "vulnerable_component": "gav://io.netty:netty-codec:4.1.66.Final",
      "impacted_artifact": "docker://elasticsearch/7.15.0-test2:latest",
      "impact_path": [
        "docker://elasticsearch/7.15.0-test2:latest",
        "generic://sha256:a53372ba228046f81171efd357179b7b02a95acccee17925e3d0295829cb42ea/sha256__a53372ba228046f81171efd357179b7b02a95acccee17925e3d0295829cb42ea.tar.gz",
        "gav://io.netty:netty-codec:4.1.66.Final"
      ],
      "path": "forks-release-local/elasticsearch/7.15.0-test2/latest/",
      "fixed_versions": [
        "4.1.68.Final"
      ],
      "published": "2021-09-12T14:10:55Z",
      "issue_id": "XRAY-184735",
      "package_type": "maven",
      "provider": "JFrog",
      "description": "Netty contains a flaw in the Bzip2Decoder class that is triggered as memory allocations for decompressed data can't be properly limited. This may allow a remote attacker to exhaust available memory resources and cause a denial of service.",
      "references": [
        "https://github.com/netty/netty/commit/41d3d61a61608f2223bb364955ab2045dd5e4020",
        "https://github.com/netty/netty/security/advisories/GHSA-grg4-wf29-r9vv",
        "http://cve.mitre.org/cgi-bin/cvename.cgi?name=2021-37136"
      ]
    }
  ]
}
`

const LicensesReportDetailsResponse = `
{
  "total_rows": 1,
  "rows" :[
      {
          "license": "MIT",
          "license_name" : "The MIT License",
          "component": "deb://debian:buster:glibc:2.28-10",
          "artifact": "docker://redis:latest-07142020122937",
          "path": "repo1/folder1/artifact",
          "artifact_scan_time": "2020-07-14T09:32:00Z",
          "unknown" : false,
          "unrecognized" : false,
          "custom" : false,
          "references": [
              "https://spdx.org/licenses/AFL-1.1.html",
              "https://spdx.org/licenses/AFL-1.1"
          ]
      }
  ]
}
`

const VulnerableXraySummaryArtifactResponse = `
{
  "artifacts": [
    {
      "general": {
        "component_id": "foo/bar:84a28a42",
        "name": "foo/bar:84a28a42",
        "path": "default/foo/bar/84a28a42/",
        "pkg_type": "Docker",
        "sha256": "c255cbe29c2da2935b4433a54e4ce6a3710490ee1d2c47bc68a7fa1732a3be24"
      },
      "issues": [
        {
          "issue_id": "XRAY-189376",
          "summary": "ImportedSymbols in debug/macho (for Open or OpenFat) in Go before 1.16.10 and 1.17.x before 1.17.3 Accesses a Memory Location After the End of a Buffer, aka an out-of-bounds slice situation.",
          "description": "ImportedSymbols in debug/macho (for Open or OpenFat) in Go before 1.16.10 and 1.17.x before 1.17.3 Accesses a Memory Location After the End of a Buffer, aka an out-of-bounds slice situation.",
          "issue_type": "security",
          "severity": "High",
          "provider": "JFrog",
          "cves": [
            {
              "cve": "CVE-2021-41771",
              "cvss_v2": "5.0/CVSS:2.0/AV:N/AC:L/Au:N/C:N/I:N/A:P",
              "cvss_v3": "7.5/CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:H",
              "cwe": [
                "CWE-119"
              ]
            }
          ],
          "created": "2021-11-09T00:00:00.702Z",
          "impact_path": [
            "default/foo/bar/84a28a42/sha256__fc66940af1388789585cf7128aeb3edc547723e307b53e59b75ad2797ac1c765.tar.gz/bar/github.com/lang/go/go"
          ],
          "components": [
            {
              "component_id": "github.com/golang/go",
              "fixed_versions": [
                "[1.16.10]",
                "[1.17.3]"
              ]
            },
            {
              "component_id": "github.com/golang/go/src",
              "fixed_versions": [
                "[1.16.10]",
                "[1.17.3]"
              ]
            }
          ],
          "component_physical_paths": [
            "sha256__fc66940af1388789585cf7128aeb3edc547723e307b53e59b75ad2797ac1c765.tar.gz/bar/github.com/lang/go/go"
          ]
        }
      ],
      "licenses": [
        {
          "components": [
            "go://github.com/golang/go:1.15.8"
          ],
          "full_name": "Unknown license",
          "more_info_url": [
            "Unknown link"
          ],
          "name": "Unknown"
        }
      ]
    }
  ]
}
`

const EntitledResponse = `
{
  "entitled": true,
  "feature_id": "contextual_analysis"
}
`

const NotEntitledResponse = `
{
  "entitled": false,
  "feature_id": "unknown"
}
`

const TriggerBuildScanResponse = `
{
  "info":"No Xray “Fail build in case of a violation” policy rule has been defined on this build. The Xray scan will run in parallel to the deployment of the build and will not obstruct the build. To review the Xray scan results, see the Xray Violations tab in the UI."
}
`

const BuildScanResultsResponse = `
{
  "build_name": "test-%[1]s",
  "build_number": "3",
  "status": "completed",
  "more_details_url": "http://localhost:8046/xray/ui/builds/test-%[1]s/3/1/xrayData?buildRepo=artifactory-build-info",
  "fail_build": false,
  "violations": [],
  "vulnerabilities": [
    {
      "cves": [
        {
          "cve": "CVE-2022-41853",
          "cvss_v3_score": "9.8",
          "cvss_v3_vector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H"
        }
      ],
      "summary": "Those using java.sql.Statement or java.sql.PreparedStatement in hsqldb (HyperSQL DataBase) to process untrusted input may be vulnerable to a remote code execution attack. By default it is allowed to call any static method of any Java class in the classpath resulting in code execution. The issue can be prevented by updating to 2.7.1 or by setting the system property \"hsqldb.method_class_names\" to classes which are allowed to be called. For example, System.setProperty(\"hsqldb.method_class_names\", \"abc\") or Java argument -Dhsqldb.method_class_names=\"abc\" can be used. From version 2.7.1 all classes by default are not accessible except those in java.lang.Math and need to be manually enabled.",
      "severity": "Critical",
      "components": {
        "gav://org.hsqldb:hsqldb:1.8.0.10": {
          "package_name": "org.hsqldb:hsqldb",
          "package_version": "1.8.0.10",
          "package_type": "maven",
          "fixed_versions": [
            "[2.7.1]"
          ],
          "infected_versions": [
            "(,2.7.1)"
          ],
          "impact_paths": [
            [
              {
                "component_id": "build://test-%[1]s:3"
              },
              {
                "component_id": "gav://org.hsqldb:hsqldb:1.8.0.10"
              }
            ]
          ]
        }
      },
      "issue_id": "XRAY-256683",
      "references": [
        "http://hsqldb.org/doc/2.0/guide/sqlroutines-chapt.html#src_jrt_access_control",
        "https://bugs.chromium.org/p/oss-fuzz/issues/detail?id=50212#c7",
        "https://lists.debian.org/debian-lts-announce/2022/12/msg00020.html"
      ],
      "is_high_profile": false,
      "provider": "JFrog",
      "edited": "0001-01-01T00:00:00Z",
      "applicability":null
    }
  ]
}
`
const xscVersionResponse = `{"xsc_version": "%s","xray_version":"3.107.8"}`

const xrayVersionResponse = `{"xray_version":"%s","xray_revision":"5735964"}`

const scanIdResponse = `{"scan_id": "3472b4e2-bddc-11ee-a9c9-acde48001122"}`

const JasConfigResponse = `{"enable_token_validation_scanning": true}`

const ArtifactStatusResponse = `{
  "overall": {
    "status": "DONE",
    "time": "2023-12-01T10:00:00Z"
  },
  "details": {
    "sca": {
      "status": "DONE",
      "time": "2023-12-01T10:00:00Z"
    },
    "contextual_analysis": {
      "status": "DONE",
      "time": "2023-12-01T10:00:00Z"
    },
    "exposures": {
      "status": "NOT_SUPPORTED",
      "time": "2023-12-01T10:00:00Z"
    },
    "violations": {
      "status": "FAILED",
      "time": "2023-12-01T10:00:00Z"
    }
  }
}`

const ArtifactStatusPendingResponse = `{
  "overall": {
    "status": "PENDING",
    "time": "2023-12-01T09:30:00Z"
  },
  "details": {
    "sca": {
      "status": "PENDING",
      "time": "2023-12-01T09:30:00Z"
    },
    "contextual_analysis": {
      "status": "NOT_SCANNED",
      "time": "2023-12-01T09:30:00Z"
    },
    "exposures": {
      "status": "NOT_SUPPORTED",
      "time": "2023-12-01T09:30:00Z"
    },
    "violations": {
      "status": "NOT_SCANNED",
      "time": "2023-12-01T09:30:00Z"
    }
  }
}`

const ArtifactStatusNotSupportedResponse = `{
  "overall": {
    "status": "NOT_SUPPORTED",
    "time": "2023-12-01T11:00:00Z"
  },
  "details": {
    "sca": {
      "status": "NOT_SUPPORTED",
      "time": "2023-12-01T11:00:00Z"
    },
    "contextual_analysis": {
      "status": "NOT_SUPPORTED",
      "time": "2023-12-01T11:00:00Z"
    },
    "exposures": {
      "status": "NOT_SUPPORTED",
      "time": "2023-12-01T11:00:00Z"
    },
    "violations": {
      "status": "NOT_SUPPORTED",
      "time": "2023-12-01T11:00:00Z"
    }
  }
}`

const XscGitInfoResponse = `{"multi_scan_id": "3472b4e2-bddc-11ee-a9c9-acde48001122"}`

const XscGitInfoBadResponse = `"failed create git info request: git_repo_url field must contain value"`

var GitInfoContextWithMinimalRequiredFields = xscServices.XscGitInfoContext{
	Source: xscServices.CommitContext{
		GitRepoHttpsCloneUrl: "https://git.jfrog.info/projects/XSC/repos/xsc-service",
		BranchName:           "feature/XRAY-123-cool-feature",
		CommitHash:           "acc5e24e69a-d3c1-4022-62eb-69e4a1e5",
	},
}

var GitInfoContextWithMissingFields = xscServices.XscGitInfoContext{
	Source: xscServices.CommitContext{
		GitRepoHttpsCloneUrl: "https://git.jfrog.info/projects/XSC/repos/xsc-service",
		BranchName:           "feature/XRAY-123-cool-feature",
	},
}

const TestMultiScanId = "3472b4e2-bddc-11ee-a9c9-acde48001122"
const TestXscVersion = "1.0.0"

var MapReportIdEndpoint = map[int]string{
	777: VulnerabilitiesEndpoint,
	888: LicensesEndpoint,
}

var MapResponse = map[string]map[string]string{
	VulnerabilitiesEndpoint: {
		"XrayReportRequest": VulnerabilityXrayReportRequestResponse,
		"ReportStatus":      VulnerabilityReportStatusResponse,
		"ReportDetails":     VulnerabilityReportDetailsResponse,
	},
	LicensesEndpoint: {
		"XrayReportRequest": LicensesXrayReportRequestResponse,
		"ReportStatus":      LicensesReportStatusResponse,
		"ReportDetails":     LicensesReportDetailsResponse,
	},
}
