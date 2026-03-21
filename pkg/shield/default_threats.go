package shield

// DefaultThreatsContent contains the built-in threat definitions.
// These provide basic protection against common attack patterns.
const DefaultThreatsContent = `# SHIELD.md - Default Threat Feed
# This file defines runtime security policies for MoonHub.

## THREAT-001 - SQL Injection
` + "```yaml" + `
id: THREAT-001
fingerprint: sql-inject-01
category: tool
severity: high
confidence: 0.90
action: block
title: SQL Injection via Tool Arguments
description: Block SQL injection attempts through tool arguments
recommendation_agent: |
  BLOCK: tool.call with arguments containing (DROP, DELETE, UNION, --, '; --, ' OR '1'='1, " OR "1"="1, ; DROP, ; DELETE, ; TRUNCATE, ; UPDATE, ; INSERT)
` + "```" + `

## THREAT-002 - Command Injection
` + "```yaml" + `
id: THREAT-002
fingerprint: cmd-inject-01
category: tool
severity: critical
confidence: 0.95
action: block
title: Command Injection via Tool Arguments
description: Block command injection patterns in tool arguments
recommendation_agent: |
  BLOCK: tool.call with arguments containing ($(` + "`" + `), $(, &&, ||, ;, |, >, <, ` + "`" + `, 2>&1, /dev/null, curl |, wget |, bash -c, sh -c, python -c, perl -e, ruby -e)
` + "```" + `

## THREAT-003 - Path Traversal
` + "```yaml" + `
id: THREAT-003
fingerprint: path-traversal-01
category: file
severity: high
confidence: 0.90
action: block
title: Path Traversal Attempt
description: Block attempts to access files outside the workspace using path traversal
recommendation_agent: |
  BLOCK: file path contains ../
  BLOCK: file path contains ..\
` + "```" + `

## THREAT-004 - Sensitive System File Access
` + "```yaml" + `
id: THREAT-004
fingerprint: sensitive-sys-01
category: file
severity: high
confidence: 0.90
action: block
title: Sensitive System File Access
description: Block access to sensitive system files
recommendation_agent: |
  BLOCK: file path contains /etc/passwd
  BLOCK: file path contains /etc/shadow
  BLOCK: file path contains /etc/sudoers
  BLOCK: file path contains /etc/ssh/
  BLOCK: file path contains ~/.ssh/
  BLOCK: file path contains .ssh/id_rsa
  BLOCK: file path contains .ssh/id_ed25519
  BLOCK: file path contains id_rsa
  BLOCK: file path contains id_ed25519
` + "```" + `

## THREAT-005 - Sensitive Credential Access
` + "```yaml" + `
id: THREAT-005
fingerprint: sensitive-cred-01
category: file
severity: high
confidence: 0.85
action: block
title: Sensitive Credential Access
description: Block access to credential files
recommendation_agent: |
  BLOCK: file path contains .env
  BLOCK: file path contains .env.local
  BLOCK: file path contains .env.production
  BLOCK: file path contains credentials.json
  BLOCK: file path contains service-account.json
  BLOCK: file path contains .gnupg/
  BLOCK: file path contains .pgp/
  BLOCK: file path contains private.key
  BLOCK: file path contains private.pem
` + "```" + `

## THREAT-006 - Prompt Injection Detection
` + "```yaml" + `
id: THREAT-006
fingerprint: prompt-inject-01
category: prompt
severity: medium
confidence: 0.75
action: log
title: Prompt Injection Pattern Detected
description: Log potential prompt injection attempts for monitoring
recommendation_agent: |
  LOG: incoming message contains ignore previous instructions
  LOG: incoming message contains ignore all previous
  LOG: incoming message contains disregard all
  LOG: incoming message contains you are now
  LOG: incoming message contains new instructions
  LOG: incoming message contains system:
  LOG: incoming message contains [SYSTEM]
  LOG: incoming message contains <system>
` + "```" + `

## THREAT-007 - Dangerous Exec Commands
` + "```yaml" + `
id: THREAT-007
fingerprint: dangerous-exec-01
category: tool
severity: critical
confidence: 0.95
action: block
title: Dangerous Exec Commands
description: Block dangerous shell commands
recommendation_agent: |
  BLOCK: tool.call exec with arguments containing (rm -rf, rm -fr, mkfs, dd if=, :(){ :|:& };:, chmod 777, chown root, curl | bash, wget | bash, nc -l, nc -e, /dev/tcp, /dev/udp)
` + "```" + `

## THREAT-008 - Network Data Exfiltration
` + "```yaml" + `
id: THREAT-008
fingerprint: data-exfil-01
category: network
severity: high
confidence: 0.85
action: log
title: Potential Data Exfiltration
description: Log suspicious outbound requests that might indicate data exfiltration
recommendation_agent: |
  LOG: outbound request to pastebin.com
  LOG: outbound request to webhook.site
  LOG: outbound request to requestbin
  LOG: outbound request to ngrok.io
  LOG: outbound request to burpcollaborator
` + "```" + `
`

// DefaultThreatIDs returns the list of default threat IDs.
func DefaultThreatIDs() []string {
	return []string{
		"THREAT-001",
		"THREAT-002",
		"THREAT-003",
		"THREAT-004",
		"THREAT-005",
		"THREAT-006",
		"THREAT-007",
		"THREAT-008",
	}
}
