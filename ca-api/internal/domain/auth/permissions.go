// pkg/auth/permissions.go
package auth

type Permission string

const (
	// Certificate Authority permissions
	CAView   Permission = "ca:view"
	CACreate Permission = "ca:create"
	CAUpdate Permission = "ca:update"
	CARotate Permission = "ca:rotate"
	CADelete Permission = "ca:delete"

	// Certificate permissions
	CertView    Permission = "cert:view"
	CertRequest Permission = "cert:request"
	CertApprove Permission = "cert:approve"
	CertRevoke  Permission = "cert:revoke"

	// Template permissions
	TemplateView   Permission = "template:view"
	TemplateCreate Permission = "template:create"
	TemplateUpdate Permission = "template:update"

	// Audit permissions
	AuditView   Permission = "audit:view"
	AuditExport Permission = "audit:export"
)

var DefaultRoles = map[string][]Permission{
	"admin": {
		CAView, CACreate, CAUpdate, CARotate, CADelete,
		CertView, CertRequest, CertApprove, CertRevoke,
		TemplateView, TemplateCreate, TemplateUpdate,
		AuditView, AuditExport,
	},
	"security_officer": {
		CAView, CARotate,
		CertView, CertApprove, CertRevoke,
		TemplateView, TemplateCreate, TemplateUpdate,
		AuditView, AuditExport,
	},
	"developer": {
		CAView,
		CertView, CertRequest,
		TemplateView,
	},
	"auditor": {
		CAView, CertView, TemplateView,
		AuditView, AuditExport,
	},
}
