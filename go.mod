module github.com/jbcom/GitWebhookProxy

// GO 1.25, UP FROM 1.13. The upstream project's last release was November 2020
// and its last commit March 2023, so it was pinned to a toolchain that had
// already left support when this fork was taken. Thirteen releases of security
// fixes, `errors.Join`, `min`/`max`, generics and a rewritten `net/http` sat on
// the other side of a one-line bump.
go 1.25

require (
	github.com/jarcoal/httpmock v1.4.1
	github.com/julienschmidt/httprouter v1.3.0
	github.com/namsral/flag v1.7.4-pre
)
