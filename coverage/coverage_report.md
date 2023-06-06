## Armstrong Test Coverage

The coverage is computed based on https://github.com/Azure/azure-rest-api-specs/tree/1c598e61adeff6930d8bce5dac374b3f8c6f1c54

<blockquote><details open><summary>/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/dataCollectionRules/{dataCollectionRuleName}?api-version=2022-06-01</summary><blockquote>

<details open><summary><span style="color:red">body(22/177)</summary><blockquote>

- <span >location</span>

- <span style="color:red">tags</span>

<details><summary><span style="color:red">identity(0/6)</span></summary><blockquote>

- <span style="color:red">userAssignedIdentities</span>

<details><summary><span style="color:red">type(0/4)</span></summary><blockquote>

- <span style="color:red">value=None</span>

- <span style="color:red">value=SystemAssigned,UserAssigned</span>

- <span style="color:red">value=SystemAssigned</span>

- <span style="color:red">value=UserAssigned</span>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">kind(0/2)</span></summary><blockquote>

- <span style="color:red">value=Linux</span>

- <span style="color:red">value=Windows</span>

</blockquote></details>

<details><summary><span style="color:red">properties(21/166)</span></summary><blockquote>

- <span style="color:red">dataCollectionEndpointId</span>

- <span style="color:red">description</span>

- <span style="color:red">streamDeclarations</span>

<details><summary><span style="color:red">dataFlows(4/13)</span></summary><blockquote>

- <span style="color:red">builtInTransform</span>

- <span style="color:red">destinations(1/2)</span>

- <span style="color:red">outputStream</span>

- <span style="color:red">transformKql</span>

<details><summary><span style="color:red">streams(3/6)</span></summary><blockquote>

- <span >value=Microsoft-Perf</span>

- <span >value=Microsoft-Syslog</span>

- <span >value=Microsoft-WindowsEvent</span>

- <span style="color:red">value=Microsoft-Event</span>

- <span style="color:red">value=Microsoft-InsightsMetrics</span>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">dataSources(15/115)</span></summary><blockquote>

<details><summary><span style="color:red">dataImports(0/5)</span></summary><blockquote>

<details><summary><span style="color:red">eventHub(0/4)</span></summary><blockquote>

- <span style="color:red">consumerGroup</span>

- <span style="color:red">name</span>

- <span style="color:red">stream</span>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">extensions(0/13)</span></summary><blockquote>

- <span style="color:red">extensionName</span>

- <span style="color:red">extensionSettings</span>

- <span style="color:red">inputDataSources(0/2)</span>

- <span style="color:red">name</span>

<details><summary><span style="color:red">streams(0/6)</span></summary><blockquote>

- <span style="color:red">value=Microsoft-Event</span>

- <span style="color:red">value=Microsoft-InsightsMetrics</span>

- <span style="color:red">value=Microsoft-Perf</span>

- <span style="color:red">value=Microsoft-Syslog</span>

- <span style="color:red">value=Microsoft-WindowsEvent</span>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">iisLogs(0/7)</span></summary><blockquote>

- <span style="color:red">logDirectories(0/2)</span>

- <span style="color:red">name</span>

- <span style="color:red">streams(0/2)</span>

</blockquote></details>

<details><summary><span style="color:red">logFiles(0/19)</span></summary><blockquote>

- <span style="color:red">filePatterns(0/2)</span>

- <span style="color:red">format</span>

- <span style="color:red">name</span>

- <span style="color:red">streams(0/2)</span>

<details><summary><span style="color:red">settings(0/11)</span></summary><blockquote>

<details><summary><span style="color:red">text(0/10)</span></summary><blockquote>

<details><summary><span style="color:red">recordStartTimestampFormat(0/9)</span></summary><blockquote>

- <span style="color:red">value=ISO 8601</span>

- <span style="color:red">value=M/D/YYYY HH:MM:SS AM/PM</span>

- <span style="color:red">value=MMM d hh:mm:ss</span>

- <span style="color:red">value=Mon DD, YYYY HH:MM:SS</span>

- <span style="color:red">value=YYYY-MM-DD HH:MM:SS</span>

- <span style="color:red">value=dd/MMM/yyyy:HH:mm:ss zzz</span>

- <span style="color:red">value=ddMMyy HH:mm:ss</span>

- <span style="color:red">value=yyMMdd HH:mm:ss</span>

- <span style="color:red">value=yyyy-MM-ddTHH:mm:ssK</span>

</blockquote></details>

</blockquote></details>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">performanceCounters(4/9)</span></summary><blockquote>

- <span >name</span>

- <span >samplingFrequencyInSeconds</span>

- <span style="color:red">counterSpecifiers(1/2)</span>

<details><summary><span style="color:red">streams(1/3)</span></summary><blockquote>

- <span >value=Microsoft-Perf</span>

- <span style="color:red">value=Microsoft-InsightsMetrics</span>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">platformTelemetry(0/5)</span></summary><blockquote>

- <span style="color:red">name</span>

- <span style="color:red">streams(0/2)</span>

</blockquote></details>

<details><summary><span style="color:red">prometheusForwarder(0/6)</span></summary><blockquote>

- <span style="color:red">labelIncludeFilter</span>

- <span style="color:red">name</span>

<details><summary><span style="color:red">streams(0/2)</span></summary><blockquote>

- <span style="color:red">value=Microsoft-PrometheusMetrics</span>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">syslog(8/37)</span></summary><blockquote>

- <span >name</span>

<details><summary><span style="color:red">facilityNames(2/22)</span></summary><blockquote>

- <span >value=cron</span>

- <span >value=syslog</span>

- <span style="color:red">value=*</span>

- <span style="color:red">value=auth</span>

- <span style="color:red">value=authpriv</span>

- <span style="color:red">value=daemon</span>

- <span style="color:red">value=kern</span>

- <span style="color:red">value=local0</span>

- <span style="color:red">value=local1</span>

- <span style="color:red">value=local2</span>

- <span style="color:red">value=local3</span>

- <span style="color:red">value=local4</span>

- <span style="color:red">value=local5</span>

- <span style="color:red">value=local6</span>

- <span style="color:red">value=local7</span>

- <span style="color:red">value=lpr</span>

- <span style="color:red">value=mail</span>

- <span style="color:red">value=mark</span>

- <span style="color:red">value=news</span>

- <span style="color:red">value=user</span>

- <span style="color:red">value=uucp</span>

</blockquote></details>

<details><summary><span style="color:red">logLevels(4/10)</span></summary><blockquote>

- <span >value=Alert</span>

- <span >value=Critical</span>

- <span >value=Debug</span>

- <span >value=Emergency</span>

- <span style="color:red">value=*</span>

- <span style="color:red">value=Error</span>

- <span style="color:red">value=Info</span>

- <span style="color:red">value=Notice</span>

- <span style="color:red">value=Warning</span>

</blockquote></details>

<details><summary><span style="color:red">streams(1/2)</span></summary><blockquote>

- <span >value=Microsoft-Syslog</span>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">windowsEventLogs(3/8)</span></summary><blockquote>

- <span >name</span>

- <span style="color:red">xPathQueries(1/2)</span>

<details><summary><span style="color:red">streams(1/3)</span></summary><blockquote>

- <span >value=Microsoft-WindowsEvent</span>

- <span style="color:red">value=Microsoft-Event</span>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">windowsFirewallLogs(0/5)</span></summary><blockquote>

- <span style="color:red">name</span>

- <span style="color:red">streams(0/2)</span>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">destinations(2/34)</span></summary><blockquote>

<details><summary><span style="color:red">azureMonitorMetrics(0/2)</span></summary><blockquote>

- <span style="color:red">name</span>

</blockquote></details>

<details><summary><span style="color:red">eventHubs(0/4)</span></summary><blockquote>

- <span style="color:red">eventHubResourceId</span>

- <span style="color:red">name</span>

</blockquote></details>

<details><summary><span style="color:red">eventHubsDirect(0/4)</span></summary><blockquote>

- <span style="color:red">eventHubResourceId</span>

- <span style="color:red">name</span>

</blockquote></details>

<details><summary><span style="color:red">logAnalytics(2/4)</span></summary><blockquote>

- <span >name</span>

- <span >workspaceResourceId</span>

</blockquote></details>

<details><summary><span style="color:red">monitoringAccounts(0/4)</span></summary><blockquote>

- <span style="color:red">accountResourceId</span>

- <span style="color:red">name</span>

</blockquote></details>

<details><summary><span style="color:red">storageAccounts(0/5)</span></summary><blockquote>

- <span style="color:red">containerName</span>

- <span style="color:red">name</span>

- <span style="color:red">storageAccountResourceId</span>

</blockquote></details>

<details><summary><span style="color:red">storageBlobsDirect(0/5)</span></summary><blockquote>

- <span style="color:red">containerName</span>

- <span style="color:red">name</span>

- <span style="color:red">storageAccountResourceId</span>

</blockquote></details>

<details><summary><span style="color:red">storageTablesDirect(0/5)</span></summary><blockquote>

- <span style="color:red">name</span>

- <span style="color:red">storageAccountResourceId</span>

- <span style="color:red">tableName</span>

</blockquote></details>

</blockquote></details>

</blockquote></details>

</blockquote></details>

</blockquote></details>
</blockquote>
